import { createBrowserRouter, RouterProvider } from "react-router-dom"
import Base from "./components/layouts/Base"
import Projects from "./components/pages/Projects"
import ProjectDetails from "./components/pages/ProjectDetails"
import Jobs from "./components/pages/Jobs"
import JobDetails from "./components/pages/JobDetails"
import Volumes from "./components/pages/Volumes"
import Settings from "./components/pages/Settings"


const Router= () => {
    const router = createBrowserRouter([
        {
            path: "/",
            element: <Base/>,
            children: [
                {
                    index: true,
                    element: <Projects />
                },
                {
                    path: "/projects/:projectId",
                    element: <ProjectDetails />
                },
                {
                    path: "/jobs",
                    element: <Jobs />
                },
                {
                    path: "/jobs/:jobId",
                    element: <JobDetails />
                },
                {
                    path: "/volumes",
                    element: <Volumes />
                },
                {
                    path: "/settings",
                    element: <Settings />
                }
            ]
        },  
    ])

    return <RouterProvider router={router} />;
}

export default Router;