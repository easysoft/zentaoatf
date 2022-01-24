export interface Config {
  language: string
  url: string;
  username: string;
  password: string;

  javascript: string
  lua:        string
  perl:       string
  php:        string
  python:     string
  ruby:       string
  tcl:        string
  autoit:     string
}

export interface Interpreter {
  lang: string;
  val: string;
}
