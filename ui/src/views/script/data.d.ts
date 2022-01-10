export interface Script {
  id: number;
  name: string;
  desc: string;
  content: string
}

export interface QueryResult {
  list: Script[];
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


export type Script = {
  content: string,
}

export type DeepTestMsg = {
  scope: string
  content: DeepTestMsgContent
}
export type DeepTestMsgContent = {
  act: string
  mainWindowId: number
  recorderWindowId: number
  recorderTabId: number

  data: DeepTestMsgOpt
}
export type DeepTestMsgOpt = {
  selector: any
  value: any
  tagName: any
  action: string
  keyCode: number
  href: string
  coordinates: any
}