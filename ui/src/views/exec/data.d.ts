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
