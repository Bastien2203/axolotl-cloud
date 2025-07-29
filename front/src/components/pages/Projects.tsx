import { useEffect, useState } from "react";

import { createProject, getProjects } from "../../api/projects";
import Button from "../atoms/Button";
import { Plus } from "lucide-react";
import CreateProjectModal from "../modals/CreateProjectModal";
import { useToast } from "../../contexts/ToastContext";
import type { Project } from "../../api/types";
import { useNavigate } from "react-router-dom";


const Projects = () => {
    const [projects, setProjects] = useState<Project[]>([])
    const [createModalOpen, setCreateModalOpen] = useState(false);
    const toast = useToast();
    const navigate = useNavigate();

    useEffect(() => {
        getProjects().then(setProjects)
    }, [])

    const handleProjectClick = (project: Project) => {
        navigate(`/projects/${project.id}`);
    }


    return (
        <div className="w-full m-4">
            {createModalOpen && <CreateProjectModal onClose={() => setCreateModalOpen(false)} onCreate={(project) => {
                createProject(project).then((newProject) => {
                    setProjects((prev) => [...prev, newProject]);
                    toast.success("Project created successfully!");
                }).catch((error) => {
                    console.error("Failed to create project:", error);
                    toast.error("Failed to create project. Please try again.");
                }).finally(() => {
                    setCreateModalOpen(false);
                });
            }} />}

            <div className="flex justify-end items-center w-full">
                <Button onClick={() => setCreateModalOpen(true)}>
                    Create Project <Plus />
                </Button>

            </div>
            <div className="grid grid-cols-[repeat(auto-fill,minmax(10em,1fr))] gap-4 p-4">
            {projects.map(p => (
                <ProjectCard key={p.id} project={p} onClick={() => handleProjectClick(p)} />
            ))}
            </div>
        </div>
    );
}


const ProjectCard = ({ project, onClick }: { project: Project, onClick: () => void }) => {
    return (
        <div className="shadow p-4 rounded w-[10em] aspect-square flex flex-col items-center justify-between hover:opacity-80 hover:bg-gray-100 cursor-pointer relative bg-white" onClick={onClick}>
            <div></div>
            <img src={project.icon_url} alt={`${project.name} icon`} className="w-12 h-auto" />
            <h3 >{project.name}</h3>
        </div>
    );
}



export default Projects;