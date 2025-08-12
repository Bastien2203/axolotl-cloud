import { useEffect, useState } from "react";

import { createProject, deleteProject, getProjects, updateProject } from "../../api/projects";
import Button from "../atoms/Button";
import { Plus } from "lucide-react";
import { useToast } from "../../contexts/ToastContext";
import type { Project } from "../../api/types";
import { useNavigate } from "react-router-dom";
import ProjectCard from "../all/ProjectCard";
import FormModal from "../modals/FormModal";
import { useDialog } from "../../hooks/useDialog";


const Projects = () => {
    const [projects, setProjects] = useState<Project[]>([])
    const toast = useToast();
    const navigate = useNavigate();
    const {dialog, openDialog, closeDialog} = useDialog<"create-project">();

    useEffect(() => {
        getProjects().then(setProjects)
    }, [])

    const handleProjectClick = (project: Project) => {
        navigate(`/projects/${project.id}`);
    }

    const handleDeleteProject = (project: Project) => {
        deleteProject(project.id).then(() => {
            setProjects((prev) => prev.filter(p => p.id !== project.id));
            toast.success("Project deleted successfully!");
        }).catch((error) => {
            console.error("Failed to delete project:", error);
            toast.error("Failed to delete project. Please try again.");
        });
    }

    const handleEdit = (project: Project) => {
        updateProject(project.id, project).then((updatedProject) => {
            setProjects((prev) => prev.map(p => p.id === updatedProject.id ? updatedProject : p));
            toast.success("Project updated successfully!");
        }).catch((error) => {
            console.error("Failed to update project:", error);
            toast.error("Failed to update project. Please try again.");
        });
    };

    const handleCreate = (project: Project) => {
        createProject(project).then((newProject) => {
            setProjects((prev) => [...prev, newProject]);
            toast.success("Project created successfully!");
            closeDialog("create-project");
        }).catch((error) => {
            console.error("Failed to create project:", error);
            toast.error("Failed to create project. Please try again.");
        });
    };


    return (
        <>
            <h1 className="text-2xl font-bold mb-4">Projects</h1>
            {
                dialog("create-project", (
                    <FormModal 
                    name="Project"
                    onClose={() => closeDialog("create-project")}
                    onSubmit={handleCreate}
                    fields={[
                        { name: "name", type: "text", required: true, placeholder: "Project Name" },
                        { name: "icon_url", type: "text", required: true, placeholder: "Icon URL" },
                        { name: "website_url", type: "url", required: false, placeholder: "Website URL" }
                    ]}
                />
                ))
            }

            <div className="flex justify-end items-center w-full">
                <Button onClick={() => openDialog("create-project")}>
                    Create Project <Plus />
                </Button>

            </div>
            <div className="grid grid-cols-[repeat(auto-fill,minmax(10em,1fr))] gap-4 p-4">
                {projects.map(p => (
                    <ProjectCard key={p.id} project={p} onClick={() => handleProjectClick(p)} onDelete={handleDeleteProject} onEdit={handleEdit} />
                ))}
            </div>
        </>
    );
}




export default Projects;