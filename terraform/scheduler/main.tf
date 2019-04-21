variable "topic_name" {
  type        = "string"
  description = "cloud pub/sub topic name to publish to, when shedule comes"
}

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

resource "google_pubsub_topic" "topic" {
  name = "${var.topic_name}"
}

resource "google_cloud_scheduler_job" "job" {
  provider    = "google-beta"
  name        = "${var.scheduler_name}"
  region      = "${var.region}"
  description = "trigger the slack status app cloud function every 1 hour"
  schedule    = "${var.schedule}"
  time_zone   = "${var.timezone}"

  pubsub_target {
    topic_name = "${google_pubsub_topic.topic.id}"
    data       = "${base64encode("status")}"
  }
}
