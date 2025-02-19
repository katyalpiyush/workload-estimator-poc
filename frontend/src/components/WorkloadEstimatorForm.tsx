import React, { useState } from 'react';

interface FormData {
  no_of_documents: number;
  average_document_size: number;
  workload_nature: string;
}

interface WorkloadEstimatorFormProps {
  onSubmit: (data: FormData) => void;
}

const WorkloadEstimatorForm: React.FC<WorkloadEstimatorFormProps> = ({ onSubmit }) => {
  const [noOfDocuments, setNoOfDocuments] = useState(100000);
  const [averageDocumentSize, setAverageDocumentSize] = useState(5000);
  const [workloadNature, setWorkloadNature] = useState("read");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({
      no_of_documents: noOfDocuments,
      average_document_size: averageDocumentSize,
      workload_nature: workloadNature,
    });
  };

  const handleWorkloadNatureChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setWorkloadNature(event.target.value);
  };

  return (
    <form onSubmit={handleSubmit} className="max-w-lg mx-auto p-6 bg-white rounded-lg shadow-xl">
      <div className="mb-6">
        <label className="block text-lg font-semibold text-gray-700 mb-2">Number of Documents</label>
        <input
          type="number"
          value={noOfDocuments}
          onChange={(e) => setNoOfDocuments(Number(e.target.value))}
          className="w-full p-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

      <div className="mb-6">
        <label className="block text-lg font-semibold text-gray-700 mb-2">Average Document Size (bytes)</label>
        <input
          type="number"
          value={averageDocumentSize}
          onChange={(e) => setAverageDocumentSize(Number(e.target.value))}
          className="w-full p-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

      {/* Workload Nature Radio Cards */}
      <div className="mb-6">
        <label className="block text-lg font-semibold text-gray-700 mb-2">Workload Nature</label>
        <div className="flex space-x-4">
          <label className="flex-1 cursor-pointer">
            <input
              type="radio"
              name="workload_nature"
              value="read"
              checked={workloadNature === 'read'}
              onChange={handleWorkloadNatureChange}
              className="hidden"
            />
            <div
              className={`p-6 border-2 rounded-lg shadow-md transition-all duration-300 ease-in-out transform hover:scale-105 ${
                workloadNature === 'read' ? 'bg-blue-100 border-blue-500' : 'bg-white border-gray-300'
              }`}
            >
              <h3 className="text-xl font-bold text-center text-blue-600">Read</h3>
              <p className="text-sm text-center text-gray-600">A read workload is used for querying data.</p>
            </div>
          </label>

          <label className="flex-1 cursor-pointer">
            <input
              type="radio"
              name="workload_nature"
              value="write"
              checked={workloadNature === 'write'}
              onChange={handleWorkloadNatureChange}
              className="hidden"
            />
            <div
              className={`p-6 border-2 rounded-lg shadow-md transition-all duration-300 ease-in-out transform hover:scale-105 ${
                workloadNature === 'write' ? 'bg-green-100 border-green-500' : 'bg-white border-gray-300'
              }`}
            >
              <h3 className="text-xl font-bold text-center text-green-600">Write</h3>
              <p className="text-sm text-center text-gray-600">A write workload is used for inserting data.</p>
            </div>
          </label>

          <label className="flex-1 cursor-pointer">
            <input
              type="radio"
              name="workload_nature"
              value="readwrite"
              checked={workloadNature === 'readwrite'}
              onChange={handleWorkloadNatureChange}
              className="hidden"
            />
            <div
              className={`p-6 border-2 rounded-lg shadow-md transition-all duration-300 ease-in-out transform hover:scale-105 ${
                workloadNature === 'readwrite' ? 'bg-yellow-100 border-yellow-500' : 'bg-white border-gray-300'
              }`}
            >
              <h3 className="text-xl font-bold text-center text-yellow-600">Read/Write</h3>
              <p className="text-sm text-center text-gray-600">A read-write workload involves both querying and inserting data.</p>
            </div>
          </label>
        </div>
      </div>

      <button
        type="submit"
        className="w-full bg-blue-500 text-white p-3 rounded-lg mt-6 hover:bg-blue-600 focus:outline-none focus:ring-4 focus:ring-blue-300"
      >
        Estimate Resources
      </button>
    </form>
  );
};

export default WorkloadEstimatorForm;
