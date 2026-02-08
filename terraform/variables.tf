variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "GCP region for resources"
  type        = string
  default     = "asia-southeast1"
}

variable "zone" {
  description = "GCP zone for resources"
  type        = string
  default     = "asia-southeast1-a"
}

variable "bot_token" {
  description = "Telegram bot token"
  type        = string
  sensitive   = true
}

variable "open_ai_api_key" {
  description = "OpenAI API key"
  type        = string
  sensitive   = true
}

variable "google_credentials_json" {
  description = "Google credentials JSON content"
  type        = string
  sensitive   = true
}

