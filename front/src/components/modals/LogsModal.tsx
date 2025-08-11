import { useEffect, useState } from "react";
import Modal from "../atoms/Modal";
import { getContainerLogs } from "../../api/containers";
import { useToast } from "../../contexts/ToastContext";


const LogsModal = ({ projectId, containerId, onClose }: { projectId: string, containerId: string; onClose: () => void; }) => {
    const [logs, setLogs] = useState<string>("");
    const toast = useToast();

    useEffect(() => {
        getContainerLogs(projectId, containerId, 100).then(setLogs).catch((error) => {
            console.error("Failed to fetch logs:", error);
            toast.error("Failed to fetch logs");
        });
    }, [])

    return <Modal onClose={onClose}>
        <div className="min-w-[60vw] min-h-[40vh]">
        <h2 className="text-xl font-semibold mb-4">Container Logs</h2>
        <pre className="whitespace-pre-wrap break-words">{logs == "" ? "No logs available" : logs}</pre>
        </div>
    </Modal>
}

export default LogsModal;