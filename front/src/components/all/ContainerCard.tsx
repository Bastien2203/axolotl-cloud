import { Square, Play } from "lucide-react";
import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import clsx from "clsx";

import type {
    Container,
    ContainerStatus,
    Job,
} from "../../api/types";
import { statusColors } from "../../api/types";
import { useToast } from "../../contexts/ToastContext";
import { hideModal, showModal } from "../../libs/utils/modal";

import Spinner from "../atoms/Spinner";
import JobStatusIcon from "../atoms/JobStatus";
import MoreMenu from "../atoms/MoreMenu";
import CreateContainerModal from "../modals/CreateContainerModal";
import ValidationModal from "../modals/ValidationModal";
import { getJob } from "../../api/jobs";

const ContainerCard = ({
    container,
    onDelete,
    onEdit,
    containerStatus,
    startContainer,
    stopContainer,
}: {
    container: Container;
    onDelete: (c: Container) => void;
    onEdit: (c: Container) => void;
    containerStatus: (id: string) => Promise<ContainerStatus>;
    startContainer: (id: string) => Promise<{ job_id: string }>;
    stopContainer: (id: string) => Promise<{ job_id: string }>;
}) => {
    const [status, setStatus] = useState<ContainerStatus>();
    const [loading, setLoading] = useState(false);
    const [freshJob, setFreshJob] = useState<Job | null>(container.last_job || null);
    const toast = useToast();

    useEffect(() => {
        const fetchStatus = async () => {
            try {
                setStatus(await containerStatus(container.id));
            } catch {
                setStatus("dead");
            }
        };
        fetchStatus();
        const iv = setInterval(fetchStatus, 5000);
        return () => clearInterval(iv);
    }, [container.id]);

    const waitJobCompletion = async (jobId: string) => {
        let job = await getJob(jobId);
        setFreshJob(job);
        while (job.status === "running" || job.status === "pending") {
            await new Promise((r) => setTimeout(r, 2000));
            job = await getJob(jobId);
            setFreshJob(job);
        }
        return job;
    };

    const handleToggle = async () => {
        setLoading(true);
        setStatus("loading");
        try {
            const currentStatus = status;
            const { job_id } = await (
                currentStatus === "running" ? stopContainer(container.id) : startContainer(container.id)
            );

            const finalJob = await waitJobCompletion(job_id);

            const newStatus = await containerStatus(container.id);
            setStatus(newStatus);

            if (finalJob.status === "completed") {
                toast.success(
                    status === "running"
                        ? "Container stopped successfully!"
                        : "Container started successfully!"
                );
            } else {
                toast.error(
                    status === "running"
                        ? "Container failed to stop."
                        : "Container failed to start."
                );
            }
        } catch (e) {
            console.error(e);
            toast.error(`Error ${status === "running" ? "stopping" : "starting"} container.`);
        } finally {
            setLoading(false);
        }
    };

    const badgeClass = clsx(
        "px-2 py-1 text-xs rounded-full font-medium",
        status ? statusColors[status] : "bg-gray-300"
    );

    return (
        <>
            <dialog id={`edit-container-modal-${container.id}`}>
                <CreateContainerModal
                    defaultValue={container}
                    onClose={() => hideModal(`edit-container-modal-${container.id}`)}
                    onCreate={(data) => {
                        onEdit({ ...container, ...data });
                        hideModal(`edit-container-modal-${container.id}`);
                    }}
                />
            </dialog>

            <dialog id={`validation-modal-container-delete-${container.id}`}>
                <ValidationModal
                    text="Are you sure you want to delete this container?"
                    label="Delete"
                    variant="danger"
                    onClose={() =>
                        hideModal(`validation-modal-container-delete-${container.id}`)
                    }
                    onConfirm={() => {
                        onDelete(container);
                        hideModal(`validation-modal-container-delete-${container.id}`);
                    }}
                />
            </dialog>

            <MoreMenu
                absolute
                options={[
                    {
                        label: "Edit",
                        onClick: () => showModal(`edit-container-modal-${container.id}`),
                    },
                    {
                        label: "Delete",
                        variant: "danger",
                        onClick: () => showModal(`validation-modal-container-delete-${container.id}`),
                    },
                ]}
            >
               <div className="relative p-6 rounded-2xl shadow-sm bg-gradient-to-br from-white to-gray-50 border border-gray-200 hover:shadow-md transition-all space-y-5 group">
    <div className="flex justify-between items-start">
        <div className="space-y-1">
            <div className="flex items-center gap-2">
                <h3 className="text-lg font-semibold text-gray-800 group-hover:text-gray-900">{container.name}</h3>
                <span className={badgeClass}>{status ?? "Unknown"}</span>
            </div>
            <p className="text-sm text-gray-500">{container.docker_image}</p>
        </div>
    </div>

    <div className="grid gap-3 text-sm">
        {[
            { label: "Ports", data: container.ports },
            { label: "Environment", data: container.env },
            { label: "Volumes", data: container.volumes },
        ].map(({ label, data }) => (
            <details key={label} className="bg-white/50 border border-gray-100 rounded-md p-3 hover:bg-white transition">
                <summary className="cursor-pointer font-medium text-gray-700">{label}</summary>
                <ul className="ml-4 list-disc mt-1 text-gray-600">
                    {data && Object.entries(data).length > 0 ? (
                        Object.entries(data).map(([k, v]) => <li key={k}><span className="font-medium">{k}</span> → {v}</li>)
                    ) : (
                        <li className="italic text-gray-400">None</li>
                    )}
                </ul>
            </details>
        ))}
    </div>

    {freshJob && (
        <Link
            to={`/jobs/${freshJob.id}`}
            className="text-sm text-blue-600 hover:underline flex items-center gap-1"
        >
            Last Job: {freshJob.name} <JobStatusIcon status={freshJob.status} />
        </Link>
    )}

    <div className="flex justify-end">
        <button
            onClick={handleToggle}
            disabled={loading}
            className={clsx(
                "flex items-center gap-2 text-sm font-medium rounded px-3 py-1.5 transition border",
                status === "running"
                    ? "text-red-700 border-red-200 hover:bg-red-50"
                    : "text-blue-700 border-blue-200 hover:bg-blue-50",
                loading && "opacity-50 cursor-not-allowed"
            )}
        >
            {loading ? (
                <>
                    <Spinner size={16} /> Loading
                </>
            ) : status === "running" ? (
                <>
                    <Square size={16} /> Stop
                </>
            ) : (
                <>
                    <Play size={16} /> Start
                </>
            )}
        </button>
    </div>
</div>

            </MoreMenu>
        </>
    );
};

export default ContainerCard;
