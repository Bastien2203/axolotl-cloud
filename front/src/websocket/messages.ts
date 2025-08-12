import type { JobLog } from "../api/types";


export type WSMessage = JobLogUpdateMessage
  | SubscribeMessage
  | UnsubscribeMessage;

export type SubscribeMessage = {
    type: 'subscribe';
    data: string;
}

export type UnsubscribeMessage = {
    type: 'unsubscribe';
    data: string;
}

export type JobLogUpdateMessage = {
    type: 'job_log_update';
    data: {
        jobId: string;
        log: JobLog;
    };
}