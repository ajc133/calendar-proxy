use anyhow::Result;
use calendar_proxy::CalendarParams;
use reqwest::header;

#[tokio::main]
async fn main() -> Result<()> {
    tracing_subscriber::fmt::init();

    tracing::info!("Init client");
    let client = reqwest::Client::new();

    tracing::info!("Making request");
    let param_string = CalendarParams {
        url: String::from("hello"),
        replacement_summary: String::from("world"),
    }
    .to_url_form_encoded();

    println!("{}", param_string);
    let res: reqwest::Response = client
        .post("http://localhost:3000/calendar")
        .header(
            header::CONTENT_TYPE,
            "application/x-www-form-urlencoded; charset=utf-8",
        )
        .body(param_string)
        .send()
        .await?;
    println!("{}", res.text().await.unwrap());

    Ok(())
}
