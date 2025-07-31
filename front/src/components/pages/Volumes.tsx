import { useEffect, useState } from "react";
import { getVolumes } from "../../api/volumes";
import { useToast } from "../../contexts/ToastContext";
import { type Volume } from "../../api/types";


const Volumes = () => {
    const [volumes, setVolumes] = useState<Volume[]>();
    const toast = useToast();

    useEffect(() => {
        getVolumes()
            .then(setVolumes)
            .catch((error) => {
                console.error("Failed to fetch volumes:", error);
                toast.error("Failed to fetch volumes. Please try again later.");
            });
    }, []);

    const bytesToMB = (bytes: number) => {
        return (bytes / (1024 * 1024)).toFixed(2);
    };


    if (!volumes) {
        return <div className="text-center text-gray-500">Loading volumes...</div>;
    }

    if (volumes.length === 0) {
        return <div className="text-center text-gray-500">No volumes found.</div>;
    }

    return (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {volumes.map((volume) => (
                <div
                    key={`${volume.container_id}-${volume.project_id}-${volume.destination}`}
                    className="p-4 bg-white border border-gray-300 rounded-xl shadow-sm"
                >
                    <div className="flex justify-between items-center">
                        <h3 className="text-lg font-semibold">{volume.destination}</h3>
                        <span className="text-xs text-gray-400">{volume.type}</span>
                    </div>
                    <p className="text-sm text-gray-600 truncate">{volume.source}</p>
                    <p className="text-sm text-gray-500 mt-1">
                        Size: {bytesToMB(volume.size)} MB
                    </p>
                </div>
            ))}
        </div>
    );
}
export default Volumes;