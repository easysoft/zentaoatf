export interface Execution {
  id: number;
  name: string;
  desc: string;
}

export interface ExecutionBy {
  productId: string
  moduleId: string
  suiteId: string
  taskId: string
}

export type ExecutionItem = {
  steps: StepItem[];
}
export type StepItem = {
  action: string;
  selector: string;
  value: string;
  image: string;
}

export interface WsMsg {
  msg:       string
  isRunning: string
  category:  string
}
