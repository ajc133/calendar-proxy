use axum::{
    extract::Query,
    http::header,
    response::{IntoResponse, Response},
    Form,
};
use icalendar::{self, parser::read_calendar, parser::unfold, parser::Calendar};
use reqwest::StatusCode;
use serde::Deserialize;
use std::{env, path::Path};
use uuid::Uuid;

#[derive(Debug, Deserialize)]
pub struct IDRequest {
    id: String,
}

#[derive(Debug, Deserialize)]
pub struct CalendarParams {
    pub url: String,
    pub replacement_summary: String,
}

impl CalendarParams {
    pub fn to_url_form_encoded(&self) -> String {
        format!(
            "url={}&replacement_summary={}",
            self.url, self.replacement_summary
        )
    }
}

pub async fn handle_get_calendar(Query(id_req): Query<IDRequest>) -> impl IntoResponse {
    let id = id_req.id;
    println!("GET /calendar?id={}", id);
    let calendar_params = read_record(id).await.unwrap();
    let calendar_str = fetch_calendar_text(&calendar_params.url).await;
    let unfolded = unfold(&calendar_str);

    let mut calendar = match read_calendar(&unfolded) {
        Ok(calendar) => calendar,
        Err(err) => {
            tracing::error!("Unable to parse {}: {}", &calendar_params.url, &err);
            return Response::builder()
                .status(StatusCode::UNPROCESSABLE_ENTITY)
                .body(format!(
                    "Error parsing calendar at given url: {}",
                    &calendar_params.url
                ))
                .unwrap();
        }
    };

    replace_summary(&mut calendar, calendar_params.replacement_summary);

    Response::builder()
        .header(header::CONTENT_TYPE, "text/calendar; charset=utf-8")
        .body(calendar.to_string())
        .unwrap()
}

pub async fn handle_post_calendar(
    Form(calendar_params): Form<CalendarParams>,
) -> impl IntoResponse {
    let id = write_record(calendar_params).await.unwrap();
    Response::builder()
        .header(header::CONTENT_TYPE, "text/html; charset=utf-8")
        .body(id)
        .unwrap()
}

async fn read_record(uuid: String) -> rusqlite::Result<CalendarParams, rusqlite::Error> {
    let db_dir = env::var("DATA_DIRECTORY").unwrap_or(".".to_string());
    let db_dir = Path::new(&db_dir);
    let db_path = db_dir.join("db.sqlite");
    let connection = rusqlite::Connection::open(&db_path).unwrap();

    let mut statement = connection
        .prepare("SELECT url, replacement_summary FROM calendars WHERE id = (?1) LIMIT 1")?;
    let mut params_iter = statement
        .query_map([uuid.as_str()], |row| {
            Ok(CalendarParams {
                url: row.get(0)?,
                replacement_summary: row.get(1)?,
            })
        })
        .unwrap();
    params_iter.next().unwrap()
}

async fn write_record(
    calendar_params: CalendarParams,
) -> rusqlite::Result<String, rusqlite::Error> {
    let db_dir = env::var("DATA_DIRECTORY").unwrap_or(".".to_string());
    let db_dir = Path::new(&db_dir);
    let db_path = db_dir.join("db.sqlite");
    println!("Opening {:?}", &db_path);
    let connection = rusqlite::Connection::open(&db_path).unwrap();

    let id = Uuid::new_v4()
        .hyphenated()
        .encode_lower(&mut Uuid::encode_buffer())
        .to_string();
    println!("{}", id);
    let _ = connection.execute(
        "CREATE TABLE IF NOT EXISTS calendars (
                id TEXT PRIMARY KEY,
                url TEXT, 
                replacement_summary TEXT
            )",
        (),
    );
    connection.execute(
        "INSERT INTO calendars (id, url, replacement_summary)
            VALUES (?1, ?2, ?3);",
        (
            &id,
            calendar_params.url,
            calendar_params.replacement_summary,
        ),
    )?;
    Ok(id.clone())
}

async fn fetch_calendar_text(url: &String) -> String {
    let response_result = reqwest::get(url).await.unwrap().text().await;
    tracing::info!("Fetching {}", &url);
    match response_result {
        Ok(response) => response,
        Err(err) => {
            panic!("Error fetching the calendar: {}", err);
        }
    }
}

fn replace_summary(calendar: &mut Calendar, replacement: String) {
    for component in &mut calendar.components {
        for property in &mut component.properties {
            if property.name.to_string().eq("SUMMARY") {
                property.val = replacement.clone().into();
            }
        }
    }
}
