// src/App.tsx
import React from "react";
import WorkloadEstimatorForm from "./components/WorkloadEstimatorForm";

const App: React.FC = () => {
  return (
    <div className="flex justify-center mt-10">
      <WorkloadEstimatorForm />
    </div>
  );
};

export default App;
