export interface Execution {
  id: number;
  name: string;
  desc: string;
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
