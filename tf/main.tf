resource "fly_app" "this" {
  name = "slashmovie"
}

resource "fly_ip" "this" {
  app  = fly_app.this.name
  type = "v4"
}

resource "fly_machine" "this" {
  app    = fly_app.this.name
  region = "ewr"
  image  = "chiefy/slashmovie:${var.app-version}"
  env = {
    TMDB_API_KEY         = var.tmdb-api-key
    OMDB_API_KEY         = var.omdb-api-key
    SLACK_SIGNING_SECRET = var.slack-signing-secret
  }
  services = [
    {
      ports = [
        {
          port     = 443
          handlers = ["tls", "http"]
        },
        {
          port     = 80
          handlers = ["http"]
        }
      ]
      "protocol" : "tcp",
      "internal_port" : 8080
    },
  ]
}
