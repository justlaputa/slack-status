module "scheduler" {
  source = "./scheduler"

  topic_name     = "slack-status"
  scheduler_name = "update-slack-status"
}
