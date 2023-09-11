use anyhow::Result;
use calendar_proxy::CalendarParams;
use reqwest::header;

#[tokio::main]
async fn main() -> Result<()> {
    tracing_subscriber::fmt::init();

    tracing::info!("Init client");
    let endpoint = "http://localhost:3000/calendar";
    let client = reqwest::Client::new();

    tracing::info!("Making request");
    let param_string = CalendarParams {
        url: String::from("hola"),
        replacement_summary: String::from("mundo"),
    }
    .to_url_form_encoded();
    dbg!(&param_string);

    let res: reqwest::Response = client
        .post(endpoint)
        .header(
            header::CONTENT_TYPE,
            "application/x-www-form-urlencoded; charset=utf-8",
        )
        .body(param_string)
        .send()
        .await?;

    let id: String = res.text().await.unwrap();
    dbg!(&id);

    let res = client.get(endpoint).query(&[("id", "test")]).send().await?;
    println!("{}", res.text().await.unwrap());

    Ok(())
}
