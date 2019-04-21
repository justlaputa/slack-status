provider "google" {
  credentials = "${var.credentials}"
  project     = "${var.project}"
  region      = "${var.region}"
  version     = "~> 2.5"
}

provider "google-beta" {
  credentials = "${var.credentials}"
  project     = "${var.project}"
  region      = "${var.region}"
  version     = "~> 2.5"
}
