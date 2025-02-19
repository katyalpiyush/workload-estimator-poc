import { useState } from "react";
import RecommendedConfig from "./RecommendedConfig";
import Divider from "./Divider";
import { FaCheck } from "react-icons/fa";
import { BsExclamationCircle } from "react-icons/bs";
import { FaArrowLeft } from "react-icons/fa6";

interface ServiceGroup {
  service_group_type: string;
  estimated_ram: string;
  estimated_cpu: string;
  estimated_disk: string;
  estimated_disk_io: string;
}

const WorkloadEstimatorForm: React.FC = () => {

  const [numDocuments, setNumDocuments] = useState("");
  const [docSize, setDocSize] = useState("");
  const [workloadNature, setWorkloadNature] = useState("read");
  const [serviceGroups, setServiceGroups] = useState<ServiceGroup[]>([]);

  const handleEstimate = async () => {
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
        <div>Back to &lt;start of the module&gt; </div>
      </div>
      <div className="flex flex-row gap-x-6 justify-between h-[550px]">
        <div className="flex-1 flex flex-col h-full justify-between">
          <div>
            <h1 className="text-[#181d27] text-lg font-semibold">Workload Estimator</h1>
            <p className="text-[#535861] text-sm font-normal">This estimator will give you an approximate cluster configuration to support your needs</p>
          </div>

          <Divider />

          <div className="flex flex-row gap-8">
            <div className="w-[280px]">
              <div className="text-[#414651] text-sm font-semibold">Documents</div>
              <div className="text-[#535861] text-sm font-normal">Explain what are these documents section here for user to enter right numbers</div>
            </div>

            <div className="flex flex-col justify-between gap-y-4 w-[512px]">
              <div className="flex flex-col justify-between gap-y-[6px]">
                <div className="text-[#414651] text-sm font-medium">Number of documents</div>
                <input
                  type="number"
                  placeholder="Enter number (eg: 100-1000000)"
                  className="w-full p-2 border border-[#d5d6d9] rounded-md focus:outline-none [&::-webkit-inner-spin-button]:appearance-none [&::-webkit-outer-spin-button]:appearance-none"
                  value={numDocuments}
                  onChange={(e) => setNumDocuments(e.target.value)}
                  min="0"
                />
              </div>
              <div className="flex flex-col justify-between gap-y-[6px]">
                <div className="text-[#414651] text-sm font-medium">Each document size (approx)</div>
                <input
                  type="number"
                  placeholder="Each document size (e.g. 50kb-500kb)"
                  className="w-full p-2 border border-[#d5d6d9] rounded-md focus:outline-none [&::-webkit-inner-spin-button]:appearance-none [&::-webkit-outer-spin-button]:appearance-none"
                  value={docSize}
                  onChange={(e) => setDocSize(e.target.value)}
                  min="0"
                />
              </div>
            </div>
          </div>

          <Divider />

          <div className="flex flex-row gap-8">
            <div className="w-[280px]">
              <div className="text-[#414651] text-sm font-semibold">Workload Nature</div>
              <div className="text-[#535861] text-sm font-normal">Information about this section to be here</div>
            </div>

            <div className="flex flex-col justify-between gap-y-4 w-[512px]">
              {[
                { workloadTitle: "Read", workloadSubTitle: "50% or more", description: "Text to assist this option", value: "read" },
                { workloadTitle: "Write", workloadSubTitle: "50% or more", description: "Text to assist this option", value: "write" },
                { workloadTitle: "Read and Write", workloadSubTitle: "50 - 50%", description: "Text to assist this option", value: "readwrite" },
              ].map((option) => (
                <div
                  key={option.value}
                  className={`flex flex-row gap-x-6 justify-between items-center p-4 border rounded-md cursor-pointer ${workloadNature === option.value ? "border-[#2388ff] shadow-[0_0_0_1px_#2388ff]" : "border-gray-300"
                    }`}
                  onClick={() => setWorkloadNature(option.value)}
                >
                  {/* Info Icon */}
                  <BsExclamationCircle />

                  {/* Label */}
                  <label className="flex-1 ml-3 cursor-pointer">
                    <span className="text-[#414651] font-medium"> {option.workloadTitle} </span>
                    <span className="text-[#535861] font-normal"> {option.workloadSubTitle} </span>
                    <br />
                    <span className="text-[#535861] text-sm font-normal">{option.description}</span>
                  </label>

                  {/* Square Checkbox */}
                  <div
                    className={`w-5 h-5 flex items-center justify-center border rounded-sm transition-all 
                      ${workloadNature === option.value
                        ? "bg-[#0266c2] border-[#0266c2]"
                        : "bg-white border-gray-400"
                      }`}
                  >
                    {/* Show Tick ONLY if Selected */}
                    {workloadNature === option.value && (<FaCheck className="text-[12px] text-white" />)}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        <div className="w-fit h-full">
          <RecommendedConfig serviceGroups={serviceGroups} />
        </div>
      </div>
      <div className="flex flex-row gap-x-4">
        <button onClick={handleEstimate} className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700">
          Estimate Configuration
        </button>
        <button className="text-center text-[#0266c2] text-base px-4 py-2">
          Cancel
        </button>
      </div>
    </div>
  );
};

export default WorkloadEstimatorForm;
