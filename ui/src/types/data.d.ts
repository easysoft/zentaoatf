export interface WsMsg {
  act: string;
  data: any;
}

export interface QueryResult {
  result: any[];
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
  page: number;
  pageSize: number;
  showSizeChanger: boolean;
  showQuickJumper: boolean;
}

