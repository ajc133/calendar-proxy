use axum::{routing::get, Router};
use icalendar;
use reqwest;
use serde::Deserialize;
use std::fs;
use toml;

#[derive(Deserialize)]
struct Config {
    url: String,
    replacement: String,
}

fn read_config(config_path: &str) -> Config {
    let toml_str = fs::read_to_string(config_path)
        .unwrap_or_else(|error| panic!("Failed to read '{}': {}", &config_path, error));
    toml::from_str(&toml_str)
        .unwrap_or_else(|err| panic!("Failed to parse config at {}: {}", config_path, err))
}

#[tokio::main]
async fn main() {
    let config = read_config("config.toml");

    let response_result = reqwest::get(config.url).await.unwrap().text().await;
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
                property.val = config.replacement.clone().into();
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
