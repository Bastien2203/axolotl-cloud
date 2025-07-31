import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import type { Job } from "../../api/types";
import { getJob } from "../../api/jobs";
import { useToast } from "../../contexts/ToastContext";
import { ArrowLeft } from "lucide-react";
import JobStatusIcon from "../atoms/JobStatus";

const JobDetails = () => {
    const { jobId } = useParams<{ jobId: string }>();
    const [job, setJob] = useState<Job | null>(null);
    const toast = useToast();
    const navigate = useNavigate();

    useEffect(() => {
        if (jobId) {
            getJob(jobId)
                .then(setJob)
                .catch(() => {
                    toast.error("Error fetching job details");
                });
        }
    }, [jobId]);

    if (!job) {
        return <div className="flex justify-center items-center h-screen">Loading...</div>;
    }

    return (
        <div className="p-4">
            <div className="flex items-center gap-2 mb-4">
                <ArrowLeft className="cursor-pointer" onClick={() => navigate('/jobs')} />
                <h2 className="text-xl font-semibold">{job.name}</h2>
                <i>#{job.id}</i>
                <JobStatusIcon status={job.status} />
            </div>

            {job.logs && job.logs.length > 0 ? (
                <div className="bg-gray-100 p-4 rounded-md mt-2 max-h-[70vh] overflow-auto">
                    <h3 className="text-lg font-semibold mb-2">Logs</h3>
                    <div className="text-sm text-gray-700 font-mono whitespace-pre-wrap overflow-x-auto">
                        {job.logs.map((log, index) => (
                            <div key={index}>{log.line}</div>
                        ))}
                    </div>
                </div>
            ) : (
                <div className="text-gray-500">No logs available</div>
            )}
        </div>
    );
};

export default JobDetails;
