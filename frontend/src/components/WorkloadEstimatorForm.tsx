import { useState } from "react";
import RecommendedConfig from "./RecommendedConfig";
import Divider from "./Divider";
import { FaCheck } from "react-icons/fa";
import { BsExclamationCircle } from "react-icons/bs";
import { FaArrowLeft } from "react-icons/fa6";

interface Summary {
  cluster_option: string;
  nodes_allocated: number;
  service_groups: number;
  services: string[];
  workload_type: string;
}

interface ServiceGroup {
  services: string[];
  nodes: number;
  estimated_ram: number;
  estimated_cpu: number;
  disk_type: string;
  estimated_disk: number;
  estimated_disk_io: number;
}

const WorkloadEstimatorForm: React.FC = () => {

  const [numDocuments, setNumDocuments] = useState("");
  const [docSize, setDocSize] = useState("");
  const [workloadNature, setWorkloadNature] = useState("read");
  const [summary, setSummary] = useState<Summary | null>(null);
  const [serviceGroups, setServiceGroups] = useState<ServiceGroup[]>([]);
  const [errors, setErrors] = useState({ numDocuments: false, docSize: false });

  const handleEstimate = async () => {
    const newErrors = {
      numDocuments: !numDocuments,
      docSize: !docSize,
    };

    setErrors(newErrors);

    if (newErrors.numDocuments || newErrors.docSize) {
      return; // Stop execution if errors exist
    }

    const requestBody = {
      dataset: {
        no_of_documents: Number(numDocuments),
        average_document_size: Number(docSize),
      },
      workload_nature: workloadNature,
    };

    try {
      const response = await fetch("http://localhost:8080/estimate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(requestBody),
      });

      if (response.ok) {
        const result = await response.json();
        setSummary(result.summary);
        setServiceGroups(result.service_groups_results);
      } else {
        console.error("Error fetching data:", response.statusText);
      }
    } catch (error) {
      console.error("Network error:", error);
    }
  };

  return (
    <div className="flex flex-col gap-y-8">
      <div className="flex flex-row items-center text-center text-[#0266c2] text-base font-bold gap-x-1 cursor-pointer">
        <FaArrowLeft />
        <div>Back to Create Cluster </div>
      </div>
      <div className="flex flex-row gap-x-6 justify-between">
        <div className="flex-1 flex flex-col justify-between h-[650px]">
          <div>
            <h1 className="text-[#181d27] text-lg font-semibold">Workload Estimator</h1>
            <p className="text-[#535861] text-sm font-normal">This estimator will give you an approximate cluster configuration to support your needs</p>
          </div>

          <Divider />

          <div className="flex flex-row gap-8">
            <div className="w-[280px] flex flex-col gap-y-1">
              <div className="text-[#414651] text-sm font-semibold">Documents</div>
              <div className="text-[#535861] text-sm font-normal">Please provide the document count and size to ensure accurate resource estimation.</div>
            </div>

            <div className="flex flex-col justify-between gap-y-4 w-[512px]">
              <div className="flex flex-col justify-between gap-y-[6px]">
                <div className="text-[#414651] text-sm font-medium">Number of documents</div>
                <input
                  type="number"
                  placeholder="Enter number (e.g. 100 - 1000000)"
                  className={`w-full p-2 border rounded-md focus:outline-none [&::-webkit-inner-spin-button]:appearance-none [&::-webkit-outer-spin-button]:appearance-none ${errors.numDocuments ? "border-[#ff4545] shadow-[0_0_0_1px_#ff4545]" : "border-gray-300"}`}
                  value={numDocuments}
                  onChange={(e) => setNumDocuments(e.target.value)}
                  min="0"
                />
                {errors.numDocuments && <p className="text-red-500 text-sm">This field is required.</p>}
              </div>

              <div className="flex flex-col justify-between gap-y-[6px]">
                <div className="text-[#414651] text-sm font-medium">Each document size (approx)</div>
                <input
                  type="number"
                  placeholder="Each document size (e.g. 50KB - 500KB)"
                  className={`w-full p-2 border rounded-md focus:outline-none [&::-webkit-inner-spin-button]:appearance-none [&::-webkit-outer-spin-button]:appearance-none ${errors.docSize ? "border-[#ff4545] shadow-[0_0_0_1px_#ff4545]" : "border-gray-300"}`}
                  value={docSize}
                  onChange={(e) => setDocSize(e.target.value)}
                  min="0"
                />
                {errors.docSize && <p className="text-red-500 text-sm">This field is required.</p>}
              </div>

            </div>
          </div>

          <Divider />

          <div className="flex flex-row gap-8">
            <div className="w-[280px] flex flex-col gap-y-1">
              <div className="text-[#414651] text-sm font-semibold">Workload Nature</div>
              <div className="text-[#535861] text-sm font-normal">Select the workload type to tailor resource allocation based on read and write operations.</div>
            </div>

            <div className="flex flex-col justify-between gap-y-4 w-[512px]">
              {[
                { workloadTitle: "Read", workloadSubTitle: "50% or more", description: "Optimized for read-heavy workloads", detailedDescription: "This workload primarily consists of read operations, making it suitable for read-heavy use cases.", value: "read" },
                { workloadTitle: "Write", workloadSubTitle: "50% or more", description: "Ideal for frequent data writes and updates", detailedDescription: "This workload consists mainly of write operations, ensuring data is frequently updated or inserted.", value: "write" },
                { workloadTitle: "Read and Write", workloadSubTitle: "50 - 50%", description: "Balances read and write operations equally", detailedDescription: "This workload has an equal distribution of read and write operations, balancing performance for both.", value: "readwrite" },
              ].map((option) => (
                <div
                  key={option.value}
                  className={`flex flex-row gap-x-6 justify-between items-center p-4 border rounded-md cursor-pointer relative ${workloadNature === option.value ? "border-[#2388ff] shadow-[0_0_0_1px_#2388ff]" : "border-gray-300"}`}
                  onClick={() => setWorkloadNature(option.value)}
                >
                  <div className="relative flex items-center group">
                    <BsExclamationCircle className="text-gray-500 cursor-pointer" />

                    {/* Tooltip (Fixed Positioning) */}
                    <div className="absolute left-8 top-1/2 w-[200px] max-w-xs p-2 bg-gray-800 text-white text-xs rounded shadow-md opacity-0 group-hover:opacity-100 transition-opacity duration-300 pointer-events-none break-words z-50">
                      {option.detailedDescription}
                    </div>
                  </div>

                  {/* Label */}
                  <label className="flex-1 ml-3 cursor-pointer">
                    <span className="text-[#414651] font-medium"> {option.workloadTitle} </span>
                    <span className="text-[#535861] font-normal"> {option.workloadSubTitle} </span>
                    <br />
                    <span className="text-[#535861] text-sm font-normal">{option.description}</span>
                  </label>

                  {/* Square Checkbox */}
                  <div className={`w-5 h-5 flex items-center justify-center border rounded-sm transition-all ${workloadNature === option.value ? "bg-[#0266c2] border-[#0266c2]" : "bg-white border-gray-400"}`}>
                    {workloadNature === option.value && <FaCheck className="text-[12px] text-white" />}
                  </div>
                </div>
              ))}
            </div>
          </div>

          <div className="flex flex-row gap-x-4">
            <button onClick={handleEstimate} className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 cursor-pointer">
              Estimate Configuration
            </button>
            <button
              className="text-center text-[#0266c2] text-base px-4 py-2 cursor-pointer"
              onClick={() => {
                setNumDocuments("");
                setDocSize("");
                setWorkloadNature("read");
                setErrors({ numDocuments: false, docSize: false });
                setSummary(null);
                setServiceGroups([]);
              }}
            >
              Cancel
            </button>
          </div>
        </div>

        <div className="w-fit h-full">
          <RecommendedConfig summary={summary} serviceGroups={serviceGroups} />
        </div>
      </div>
    </div>
  );
};

export default WorkloadEstimatorForm;
