variable "version" {
  type        = string
  description = "version of slashmovie image to use"
}

variable "tmdb-api-key" {
  type        = string
  sensitive   = true
  description = "tmdb API key for slashmovie app"
}

variable "omdb-api-key" {
  type        = string
  sensitive   = true
  description = "omdb API key for slashmovie app"
}

variable "slack-signing-secret" {
  type        = string
  sensitive   = true
  description = "secret token for signing slack requests"
}
