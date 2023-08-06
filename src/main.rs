use axum::{
    extract::Query,
    http::header::{self, HeaderMap},
    response::IntoResponse,
    routing::get,
    Router,
};
use icalendar::{self, parser::read_calendar, parser::unfold, parser::Calendar};
use serde::Deserialize;
use std::net::SocketAddr;

#[derive(Deserialize)]
struct CalendarParams {
    url: String,
    replacement_summary: String,
}

async fn handle_calendar(Query(calendar_params): Query<CalendarParams>) -> impl IntoResponse {
    let calendar_str = fetch_calendar_text(calendar_params.url).await;
    let unfolded = unfold(&calendar_str);

    // FIXME: Dont' panic, return 500 or something
    let mut calendar = read_calendar(&unfolded).unwrap_or_else(|err| {
        panic!("Error parsing calendar: {}", err);
    });

    replace_summary(&mut calendar, calendar_params.replacement_summary);

    let mut headers = HeaderMap::new();
    headers.insert(
        header::CONTENT_TYPE,
        "text/calendar; charset=utf-8".parse().unwrap(),
    );
    (headers, calendar.to_string())
}

async fn fetch_calendar_text(url: String) -> String {
    let response_result = reqwest::get(&url).await.unwrap().text().await;
    tracing::info!("Parsed url:{}", &url);
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
