use axum::{routing::get, Router};
use icalendar;
use reqwest;
use std::fs;
use toml::Table;

fn read_config(config_file: &str) -> String {
    let toml_str = fs::read_to_string(config_file)
        .unwrap_or_else(|_| panic!("Failed to read '{}'", &config_file));
    let calendar_toml = toml_str.parse::<Table>().unwrap();
    let Some(url) = calendar_toml["url"].as_str() else {
        panic!("Can't parse 'url' key in config.toml");
    };
    url.to_string()
}

#[tokio::main]
async fn main() {
    let url = read_config("config.toml");

    let response_result = reqwest::get(url).await.unwrap().text().await;
    let response_str = match response_result {
        Ok(response) => response,
        Err(err) => {
            panic!("Error fetching the calendar: {}", err);
        }
    };

    let mut calendar = icalendar::parser::read_calendar(&response_str).unwrap();

    for component in &mut calendar.components {
        for property in &mut component.properties {
            if property.name.to_string().eq("SUMMARY") {
                property.val = "on call".into();
            }
        }
    }

    let s = calendar.clone().to_string();
    let app = Router::new().route("/", get(|| async { s }));
    // run it with hyper on localhost:3000
    axum::Server::bind(&"0.0.0.0:3000".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
