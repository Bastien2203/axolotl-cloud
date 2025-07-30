import { useNavigate, useParams } from "react-router-dom";
import { statusColors, type Container, type ContainerStatus, type Project } from "../../api/types";
import { useEffect, useState } from "react";
import { useToast } from "../../contexts/ToastContext";
import { getProject } from "../../api/projects";
import { ArrowLeft, File, MoreVertical, Play, Plus, Square } from "lucide-react";
import { createContainer, deleteContainer, getContainers, getContainerStatus, importComposeFile, startContainer, stopContainer, updateContainer } from "../../api/containers";
import Button from "../atoms/Button";
import CreateContainerModal from "../modals/CreateContainerModal";
import ValidationModal from "../modals/ValidationModal";
import { hideModal, showModal } from "../../libs/utils/modal";
import ImportComposeFileModal from "../modals/ImportComposeFileModal";
import Spinner from "../atoms/Spinner";


const ProjectDetails = () => {
    const { projectId } = useParams<{ projectId: string }>();
    const [project, setProject] = useState<Project>();
    const [containers, setContainers] = useState<Container[]>([]);
    const [createModalOpen, setCreateModalOpen] = useState(false);
    const [importComposeFileModal, setImportComposeFileModal] = useState(false);
    const toast = useToast();
    const navigate = useNavigate();

    useEffect(() => {
        if (projectId) {
            getProject(projectId).then(setProject).catch((error) => {
                console.error("Failed to fetch project details:", error);
                toast.error("Failed to fetch project details");
            });
            getContainers(projectId)
                .then(setContainers)
                .catch((error) => {
                    console.error("Failed to fetch containers:", error);
                    toast.error("Failed to fetch containers");
                });
        }
    }, [projectId]);

    const handleDelete = (container: Container) => {
        if (!projectId) return;
        deleteContainer(projectId, container.id).then(() => {
            setContainers((prev) => prev.filter(c => c.id !== container.id));
            toast.success("Container deleted successfully!");
        }).catch((error) => {
            console.error("Failed to delete container:", error);
            toast.error("Failed to delete container. Please try again.");
        });
    }

    const handleEdit = (container: Container) => {
        if (!projectId) return;
        updateContainer(projectId, container.id, container).then(() => {
            setContainers((prev) => prev.map(c => c.id === container.id ? container : c));
            toast.success("Container updated successfully!");
        }).catch((error) => {
            console.error("Failed to update container:", error);
            toast.error("Failed to update container. Please try again.");
        });
    }


    if (!project || !projectId || !containers) {
        return <div className="w-full m-4">Loading...</div>;
    }


    return (
        <>
            {createModalOpen && <CreateContainerModal onClose={() => setCreateModalOpen(false)} onCreate={(container) => {
                createContainer(projectId, container).then((newContainer) => {
                    setContainers((prev) => [...prev, newContainer]);
                    toast.success("Container created successfully!");
                }).catch((error) => {
                    console.error("Failed to create container:", error);
                    toast.error("Failed to create container. Please try again.");
                }).finally(() => {
                    setCreateModalOpen(false);
                });
            }} />}

            {importComposeFileModal && <ImportComposeFileModal onClose={() => setImportComposeFileModal(false)} onImport={(file) => {
                importComposeFile(projectId, file).then((newContainers) => {
                    setImportComposeFileModal(false);
                    setContainers((prev) => [...prev, ...newContainers]);
                    toast.success("Containers imported successfully!");
                }).catch((error) => {
                    console.error("Failed to import containers:", error);
                    toast.error("Failed to import containers. Please try again.");
                })
            }} />}
            <div className="w-full m-4">
                <div className="flex justify-between items-center w-full">
                    <div className="flex items-center gap-2">
                        <ArrowLeft className="cursor-pointer" onClick={() => navigate('/')} />
                        <h2 className="">{project.name}</h2>
                        <img src={project.icon_url} alt={`${project.name} icon`} className="w-8 h-8 rounded-full" />
                    </div>
                    <div className="flex items-center gap-2">
                        <Button onClick={() => setImportComposeFileModal(true)} variant="secondary">
                            Import Compose File <File />
                        </Button>
                        <Button onClick={() => setCreateModalOpen(true)}>
                            Create Container <Plus />
                        </Button>
                    </div>
                </div>


                <div className="grid grid-cols-1 gap-4 mt-4">
                    {
                        containers.map((container, i) => (
                            <ContainerCard
                                key={container.id || i}
                                container={container}
                                onDelete={handleDelete}
                                onEdit={handleEdit}
                                containerStatus={(containerId) => getContainerStatus(projectId, containerId)}
                                startContainer={(containerId) => startContainer(projectId, containerId)}
                                stopContainer={(containerId) => stopContainer(projectId, containerId)}
                            />
                        ))
                    }
                </div>
            </div>
        </>
    );
}



const ContainerCard = ({ container, onDelete, onEdit, containerStatus, startContainer, stopContainer }: { container: Container, onDelete: (container: Container) => void, onEdit: (container: Container) => void, containerStatus: (id: string) => Promise<ContainerStatus>, startContainer: (id: string) => Promise<void>, stopContainer: (id: string) => Promise<void> }) => {
    const [menuOpen, setMenuOpen] = useState(false);
    const [status, setStatus] = useState<ContainerStatus>();
    const [loading, setLoading] = useState(false);
    const toggleMenu = () => setMenuOpen((prev) => !prev);
    const toast = useToast();


    const refreshStatus = () => {
        containerStatus(container.id).then(setStatus).catch((error) => {
            console.error("Failed to fetch container status:", error);
            setStatus("dead");
        });
    }

    useEffect(() => {
        const interval = setInterval(() => {
            refreshStatus();
        }, 5000);
        refreshStatus();
        return () => clearInterval(interval);
    }, [container.id]);

    const handleEdit = (updatedData: Omit<Container, 'id'>) => {
        onEdit({ ...container, ...updatedData });
        hideModal(`edit-container-modal-${container.id}`);
    };

    const handleDelete = () => {
        onDelete(container);
        hideModal(`validation-modal-container-delete-${container.id}`);
    };

    const handleToggleStatus = () => {
        setLoading(true);
        if (status === "running") {
            stopContainer(container.id).then(() => {
                setStatus("exited");
                toast.success("Container stopped successfully!");
            }).catch((error) => {
                console.error("Failed to stop container:", error);
                toast.error("Failed to stop container. Please try again.");
            }).finally(() => {
                setLoading(false);
            });
        } else {
            startContainer(container.id).then(() => {
                setStatus("running");
                toast.success("Container started successfully!");
            }).catch((error) => {
                console.error("Failed to start container:", error);
                toast.error("Failed to start container. Please try again.");
            }).finally(() => {
                setLoading(false);
            });
        }
    };

    return (
        <>
            <dialog id={`validation-modal-container-delete-${container.id}`} open={false}>
                <ValidationModal onClose={() => {
                    hideModal(`validation-modal-container-delete-${container.id}`);
                }} variant="danger" text="Are you sure you want to delete this container?" label="Delete" onConfirm={handleDelete} />
            </dialog>
            <dialog id={`edit-container-modal-${container.id}`} open={false}>
                <CreateContainerModal onClose={() => hideModal(`edit-container-modal-${container.id}`)} onCreate={handleEdit} defaultValue={container} />
            </dialog>

            <div className="relative border border-gray-200 p-4 rounded-2xl shadow-sm bg-white hover:shadow-md transition">
                <div className="flex justify-between items-start mb-2">
                    <div>
                        <div className="flex items-center gap-2 mb-1">
                            <h3 className="text-xl font-semibold">{container.name}</h3>
                            <span className={`inline-block px-2 py-1 text-xs font-medium rounded-full ${status ? statusColors[status] : "bg-gray-300"}`}>
                                {status || "Unknown"}
                            </span>
                        </div>
                        <p className="text-sm text-gray-500">{container.docker_image}</p>
                    </div>

                    <div className="relative">
                        <button onClick={toggleMenu} className="p-1 rounded-full hover:bg-gray-100">
                            <MoreVertical size={20} />
                        </button>
                        {menuOpen && (
                            <div className="absolute right-0 mt-2 w-32 bg-white border border-gray-200 rounded-lg shadow-lg z-10">
                                <button onClick={() => showModal(`edit-container-modal-${container.id}`)} className="w-full text-left px-4 py-2 text-sm hover:bg-gray-100">
                                    Edit
                                </button>
                                <button onClick={() => showModal(`validation-modal-container-delete-${container.id}`)} className="w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50">
                                    Delete
                                </button>
                            </div>
                        )}
                    </div>
                </div>

                <div className="text-sm text-gray-700 space-y-3">
                    <details className="group">
                        <summary className="cursor-pointer font-medium text-gray-800">Ports</summary>
                        <ul className="ml-4 list-disc mt-1">
                            {container.ports && Object.entries(container.ports).length > 0 ? (
                                Object.entries(container.ports).map(([hostPort, containerPort]) => (
                                    <li key={hostPort}>{hostPort} → {containerPort}</li>
                                ))
                            ) : (
                                <li className="italic text-gray-400">None</li>
                            )}
                        </ul>
                    </details>

                    <details className="group">
                        <summary className="cursor-pointer font-medium text-gray-800">Environment Variables</summary>
                        <ul className="ml-4 list-disc mt-1">
                            {container.env && Object.entries(container.env).length > 0 ? (
                                Object.entries(container.env).map(([key, value]) => (
                                    <li key={key}><span className="font-mono">{key}</span>: {value}</li>
                                ))
                            ) : (
                                <li className="italic text-gray-400">None</li>
                            )}
                        </ul>
                    </details>

                    <details className="group">
                        <summary className="cursor-pointer font-medium text-gray-800">Volumes</summary>
                        <ul className="ml-4 list-disc mt-1">
                            {container.volumes && Object.entries(container.volumes).length > 0 ? (
                                Object.entries(container.volumes).map(([host, target]) => (
                                    <li key={host}>{host} → {target}</li>
                                ))
                            ) : (
                                <li className="italic text-gray-400">None</li>
                            )}
                        </ul>
                    </details>
                </div>

                <div className="mt-4 flex justify-end">
                    <button
                        onClick={handleToggleStatus}
                        className={`flex items-center gap-1 text-sm ${status === "running" ? "text-red-600 hover:text-red-800" : "text-blue-600 hover:text-blue-800"} transition`}
                    >
                        {loading ? (
                            <Spinner />
                        ) : status === "running" ? (
                            <>
                                Stop
                                <Square size={16} />
                            </>
                        ) : (
                            <>
                                <Play size={16} />
                                Start

                            </>
                        )}
                    </button>
                </div>
            </div>
        </>
    );
};



export default ProjectDetails;