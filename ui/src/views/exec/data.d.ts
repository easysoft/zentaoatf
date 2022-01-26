export interface Execution {
  id: number;
  name: string;
  desc: string;
}

export interface WsMsg {
  msg:       string
  isRunning: string
  category:  string
}
