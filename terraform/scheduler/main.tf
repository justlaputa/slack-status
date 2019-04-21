variable "region" {
  type        = "string"
  description = "crontab schedule region"
  default     = "asia-northeast1"
}

variable "schedule" {
  type        = "string"
  description = "crontab schedule definition"

  # default to per hour
  default = "0 * * * *"
}

variable "timezone" {
  type        = "string"
  description = "time zone used for the schedule"
  default     = "Asia/Tokyo"
}

variable "scheduler_name" {
  type        = "string"
  description = "scheduler job name"
}

variable "http_trigger_url" {
  type = "string"
  description = "gcp cloud run deploy url"
}

resource "google_cloud_scheduler_job" "job" {
  provider    = "google-beta"
  name        = "${var.scheduler_name}"
  region      = "${var.region}"
  description = "trigger the slack status app cloud function every 1 hour"
  schedule    = "${var.schedule}"
  time_zone   = "${var.timezone}"

  http_target {
    http_method = "GET"
    uri = "${var.http_trigger_url}"
  }
}
