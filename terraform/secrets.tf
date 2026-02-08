resource "google_secret_manager_secret" "bot_token" {
  secret_id = "bot-token"
  project   = var.project_id

  replication {
    auto {}
  }

  depends_on = [google_project_service.apis]
}

resource "google_secret_manager_secret_version" "bot_token" {
  secret      = google_secret_manager_secret.bot_token.id
  secret_data = var.bot_token
}

resource "google_secret_manager_secret" "open_ai_api_key" {
  secret_id = "open-ai-api-key"
  project   = var.project_id

  replication {
    auto {}
  }

  depends_on = [google_project_service.apis]
}

resource "google_secret_manager_secret_version" "open_ai_api_key" {
  secret      = google_secret_manager_secret.open_ai_api_key.id
  secret_data = var.open_ai_api_key
}

resource "google_secret_manager_secret" "google_credentials" {
  secret_id = "google-credentials"
  project   = var.project_id

  replication {
    auto {}
  }

  depends_on = [google_project_service.apis]
}

resource "google_secret_manager_secret_version" "google_credentials" {
  secret      = google_secret_manager_secret.google_credentials.id
  secret_data = var.google_credentials_json
}
