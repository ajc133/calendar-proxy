use axum::{
    extract::Query,
    http::header,
    response::{IntoResponse, Response},
    Form,
};
use reqwest::StatusCode;

use icalendar::{self, parser::read_calendar, parser::unfold, parser::Calendar};
use serde::Deserialize;

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

pub async fn handle_get_calendar(
    Query(calendar_params): Query<CalendarParams>,
) -> impl IntoResponse {
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
    println!("POST: {:?}", calendar_params);
    Response::builder()
        .header(header::CONTENT_TYPE, "text/html; charset=utf-8")
        .body(String::from("ok"))
        .unwrap()
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
