module "scheduler" {
  source = "./scheduler"

  scheduler_name = "update-slack-status"
  http_trigger_url = "${var.http_trigger_url}"
}
