
import { CircleCheck, CirclePause, CirclePlay, CircleX } from "lucide-react";
import {type JobStatus} from "../../api/types";


const statusClasses = {
        pending: "bg-yellow-100 text-yellow-800",
        running: "bg-green-100 text-green-800",
        completed: "bg-blue-100 text-blue-800",
        failed: "bg-red-100 text-red-800"
    };
    
const statusIcons = {
        pending: <CirclePause />,
        running: <CirclePlay />,
        completed: <CircleCheck />,
        failed: <CircleX />
}

const JobStatusIcon = ({ status }: { status: JobStatus }) => {
    return (
        <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-semibold aspect-square ${statusClasses[status]}`}>
            {statusIcons[status]}
        </span>
    );
}

export default JobStatusIcon;