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

