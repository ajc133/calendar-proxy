use axum::{routing::post, Router};
use calendar_proxy::{handle_get_calendar, handle_post_calendar};
use std::net::SocketAddr;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();
    let addr = SocketAddr::from(([0, 0, 0, 0], 3000));
    let app = Router::new().route(
        "/calendar",
        post(handle_post_calendar).get(handle_get_calendar),
    );

    tracing::info!("listening on {}", addr);
    axum::Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .unwrap();
}
