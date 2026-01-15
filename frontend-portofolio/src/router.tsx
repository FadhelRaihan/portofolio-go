import { BrowserRouter, Routes, Route } from "react-router-dom";
import { ProjectPage } from "./pages/Portofolio/ProjectPage";

export function AppRouter() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<ProjectPage />}/>
      </Routes>
    </BrowserRouter>
  );
}
