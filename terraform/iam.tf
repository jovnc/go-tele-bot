data "google_compute_default_service_account" "default" {
  project = var.project_id
}

resource "google_secret_manager_secret_iam_member" "bot_token_access" {
  secret_id = google_secret_manager_secret.bot_token.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${data.google_compute_default_service_account.default.email}"
}

resource "google_secret_manager_secret_iam_member" "open_ai_api_key_access" {
  secret_id = google_secret_manager_secret.open_ai_api_key.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${data.google_compute_default_service_account.default.email}"
}

resource "google_secret_manager_secret_iam_member" "google_credentials_access" {
  secret_id = google_secret_manager_secret.google_credentials.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${data.google_compute_default_service_account.default.email}"
}
