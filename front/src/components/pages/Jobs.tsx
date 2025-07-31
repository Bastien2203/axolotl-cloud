import { useEffect, useState } from "react";
import type { Job } from "../../api/types";
import { getJobs, removeJob } from "../../api/jobs";
import { useToast } from "../../contexts/ToastContext";
import MoreMenu from "../atoms/MoreMenu";
import { useNavigate } from "react-router-dom";
import JobStatusIcon from "../atoms/JobStatus";


const Jobs = () => {
    const [jobs, setJobs] = useState<Job[]>();
    const toast = useToast();
    const navigate = useNavigate();


    useEffect(() => {
        getJobs().then(setJobs).catch((error) => {
            console.error("Error fetching jobs:", error);
            toast.error("Failed to fetch jobs");
        });
    }, []);

    const handleDeleteJob = async (jobId: string) => {
        removeJob(jobId)
            .then(() => {
                setJobs((prevJobs) => prevJobs?.filter((job) => job.id !== jobId));
                toast.success("Job deleted successfully");
            })
            .catch((error) => {
                console.error("Error deleting job:", error);
                toast.error("Failed to delete job");
            });
    }

    const formatDate = (date: number) => {
        const d = new Date(date * 1000);
        return `${d.toLocaleDateString()} ${d.toLocaleTimeString()}`;
    }


    return (
        <>
            <h1 className="text-2xl font-bold mb-4">Jobs</h1>

            {jobs ? 
            jobs.length > 0 ?
            (
                <div className="">
                    {jobs.map((job) => (
                        <MoreMenu key={job.id} options={[
                            {
                                label: "Delete",
                                onClick: async () => {
                                    handleDeleteJob(job.id);
                                },
                                variant: "danger"
                            }
                        ]}  onClick={() => navigate(`/jobs/${job.id}`)} className="p-4 flex justify-between items-center border-b-gray-200 border-b last:border-b-0 hover:bg-gray-50 cursor-pointer">
                            <div className="flex items-center gap-4">
                                <JobStatusIcon status={job.status} />
                                <h2 className="text-xl font-semibold">{job.name}</h2>
                                <i>#{job.id}</i>


                                <span className="text-gray-500 text-sm">{formatDate(job.created_at)}</span>

                            </div>
                        </MoreMenu>
                    ))}
                </div>
            ) : (
                <div className="text-gray-500">No jobs found</div>
            ) : (
                <div className="text-gray-500">Loading jobs...</div>
            )}

        </>
    );
}

export default Jobs;