export interface WsMsg {
  msg:       string
  isRunning: string
  category:  string
  conn: string

  time: string
  info: any
}

export interface QueryParams {
  keywords:  string,
  enabled: string,
  page: number,
  pageSize: number,
}

export interface QueryResult {
  result: any[];
  pagination: PaginationConfig;
}

export interface PaginationConfig {
  total: number;
  page: number;
  pageSize: number;
  showSizeChanger: boolean;
  showQuickJumper: boolean;
}

