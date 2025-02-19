// src/App.tsx
import React, { useState } from 'react';
import WorkloadEstimatorForm from './components/WorkloadEstimatorForm';
import ServiceGroupsTable from './components/ServiceGroupsTable';

interface ServiceGroup {
  service_group_type: string;
  estimated_ram: string;
  estimated_cpu: string;
  estimated_disk: string;
  estimated_disk_io: string;
}

const App: React.FC = () => {
  const [serviceGroups, setServiceGroups] = useState<ServiceGroup[]>([]);

  const handleFormSubmit = async (data: { no_of_documents: number; average_document_size: number; workload_nature: string }) => {
    // Updated API call to point to the correct endpoint
    const response = await fetch('http://localhost:8080/estimate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });

    // Check if the response is OK, then parse it
    if (response.ok) {
      const result = await response.json();
      setServiceGroups(result.service_groups_results);
    } else {
      console.error('Error fetching data:', response.statusText);
    }
  };

  return (
    <div className="container mx-auto p-6">
      <h1 className="text-3xl font-semibold mb-6">Workload Estimator</h1>
      <WorkloadEstimatorForm onSubmit={handleFormSubmit} />
      {serviceGroups.length > 0 && <ServiceGroupsTable serviceGroups={serviceGroups} />}
    </div>
  );
};

export default App;
