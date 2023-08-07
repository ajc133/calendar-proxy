use axum::{
    extract::Query,
    http::header,
    response::{IntoResponse, Response},
    routing::get,
    Router,
};
use icalendar::{self, parser::read_calendar, parser::unfold, parser::Calendar};
use reqwest::StatusCode;
use serde::Deserialize;
use std::net::SocketAddr;

#[derive(Deserialize)]
struct CalendarParams {
    url: String,
    replacement_summary: String,
}

async fn handle_calendar(Query(calendar_params): Query<CalendarParams>) -> impl IntoResponse {
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
        .header(
            header::CONTENT_TYPE,
            "text/calendar; charset=utf-8".parse::<String>().unwrap(),
        )
        .body(calendar.to_string())
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

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();
    let addr = SocketAddr::from(([0, 0, 0, 0], 3000));
    let app = Router::new().route("/calendar", get(handle_calendar));

    tracing::info!("listening on {}", addr);
    axum::Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .unwrap();
}
