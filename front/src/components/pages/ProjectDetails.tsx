import { useNavigate, useParams } from "react-router-dom";
import { type Container, type Project } from "../../api/types";
import { useEffect, useState } from "react";
import { useToast } from "../../contexts/ToastContext";
import { getProject } from "../../api/projects";
import { ArrowLeft, File, Plus } from "lucide-react";
import { createContainer, deleteContainer, getContainers, getContainerStatus, importComposeFile, startContainer, stopContainer, updateContainer } from "../../api/containers";
import Button from "../atoms/Button";
import CreateContainerModal from "../modals/CreateContainerModal";
import ImportComposeFileModal from "../modals/ImportComposeFileModal";
import ContainerCard from "../all/ContainerCard";


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

            <div className="flex justify-between items-center w-full">
                <div className="flex items-center gap-2">
                    <ArrowLeft className="cursor-pointer" onClick={() => navigate('/')} />
                    <h2 className="">{project.name}</h2>
                    <img src={project.icon_url} alt={`${project.name} icon`} className="w-auto h-8 " />
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
        </>
    );
}




export default ProjectDetails;