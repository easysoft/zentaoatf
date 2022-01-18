export interface WsMsg {
  act: string;
  data: any;
}

export interface QueryResult {
  data: Execution[];
  pagination: PaginationConfig;
}

export interface QueryParams {
  keywords:  string,
  enabled: string,
  page: number,
  pageSize: number,
}

export interface PaginationConfig {
  total: number;
  current: number;
  pageSize: number;
  showSizeChanger: boolean;
  showQuickJumper: boolean;
}

