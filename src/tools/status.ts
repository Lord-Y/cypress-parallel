import { Execution } from '../api/executionsService'

const executionStatus: Array<string> = [
  'NOT_STARTED',
  'QUEUED',
  'SCHEDULED',
  'RUNNING',
  'FAILED',
  'CANCELLED',
]

class Statuses {
  tests(execution: Execution): string {
    if (
      Object.keys(execution).length > 0 &&
      executionStatus.includes(execution.execution_status)
    ) {
      return execution.execution_status
    }

    if (Object.keys(execution.result).length > 0) {
      const result: any = JSON.parse(execution.result)
      let passes = 0
      const total: number = Object.keys(
        result.results[0].suites[0].tests,
      ).length
      for (let i = 0; i < total; i++) {
        if (result.results[0].suites[0].tests[i].pass) {
          passes += 1
        }
      }
      return passes + ' passes / ' + total
    }
    return ''
  }

  global(execution: Execution): string {
    if (
      Object.keys(execution).length > 0 &&
      (executionStatus.includes(execution.execution_status))
    ) {
      return execution.execution_status
    }

    if (Object.keys(execution.result).length > 0) {
      const result: any = JSON.parse(execution.result)
      let passes = 0
      const total: number = Object.keys(
        result.results[0].suites[0].tests,
      ).length
      for (let i = 0; i < total; i++) {
        if (result.results[0].suites[0].tests[i].pass) {
          passes += 1
        }
      }
      if (passes === total) {
        return 'PASSED'
      }
      return 'FAILED'
    }
    return ''
  }
}

export default new Statuses()
