variable "project" {
  type = "string"
  description = "google project to use"
  default = "laputa"
}

variable "credentials" {
    type = "string"
    description = "google service account key"
    default = "/Users/xiaohan/Documents/credentials/slack-status-infra.json"
}

variable "region" {
    type = "string"
    description = "google region"
    default = "asia-northeast1"  
}

variable "http_trigger_url" {
  type = "string"
  description = "gcp cloud run deploy url"
  default = "https://slack-status-evaw76nphq-uc.a.run.app"  
}
