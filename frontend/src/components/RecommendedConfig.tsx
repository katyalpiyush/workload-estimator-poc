import React from "react";
import { motion } from "framer-motion";

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

interface RecommendedConfigProps {
  summary: Summary | null;
  serviceGroups: ServiceGroup[];
}

const workloadLabels: Record<string, string> = {
  read: "Read (50% or more)",
  write: "Write (50% or more)",
  readwrite: "Read Write (50 - 50%)",
};

const RecommendedConfig: React.FC<RecommendedConfigProps> = ({ summary, serviceGroups }) => {
  return (
    <div className="w-[303px] h-full bg-[#F7FAFF] border border-[#E5E5E5] rounded-[12px] flex flex-col gap-y-4">
      <div className="flex items-center border-b border-[#DFDFDF] pl-4 pr-4 pb-1 mt-[5px]">
        <h2 className="text-[16px] text-[#181D27] h-[28px] flex-1 flex items-center font-medium">Recommended Configuration</h2>
      </div>

      <div className="flex-1 pl-4 pr-4">
        {serviceGroups.length === 0 ? (
          <p className="text-[#535862] text-[14px] w-[257px]">
            Start entering the information to get your recommendation here.
          </p>
        ) : (
          <motion.div
            initial={{ opacity: 0, y: 0 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.3, ease: "easeInOut" }}
          >
            <div className="flex flex-col gap-y-4">
              {summary && (
                <div className="flex flex-col gap-y-1">
                  <div className="text-neutral-800 text-[16px] font-medium">Summary</div>
                  <div>
                    <div className="flex flex-row gap-x-2"><div className="text-neutral-500 text-sm font-normal">Option: </div><div className="text-neutral-500 text-sm font-medium">{summary.cluster_option}</div></div>
                    <div className="flex flex-row gap-x-2"><div className="text-neutral-500 text-sm font-normal">Nodes Allocated: </div><div className="text-neutral-500 text-sm font-medium">{summary.nodes_allocated}</div></div>
                    <div className="flex flex-row gap-x-2"><div className="text-neutral-500 text-sm font-normal">Service Groups: </div><div className="text-neutral-500 text-sm font-medium">{summary.service_groups}</div></div>
                    <div className="flex flex-row gap-x-2"><div className="text-neutral-500 text-sm font-normal">Services: </div><div className="text-neutral-500 text-sm font-medium">{summary.services.map(service => service.charAt(0).toUpperCase() + service.slice(1)).join(", ")}</div></div>
                    <div className="flex flex-row gap-x-2"><div className="text-neutral-500 text-sm font-normal">Workload&nbsp;Type: </div><div className="text-neutral-500 text-sm font-medium">{workloadLabels[summary.workload_type] || summary.workload_type}</div></div>
                  </div>
                </div>
              )}
              {serviceGroups.length > 0 && (
                <div className="flex flex-col gap-y-2">
                  <h3 className="text-neutral-800 text-[16px] font-medium">Detailed Breakdown</h3>
                  {serviceGroups.map((group, index) => (
                    <div key={index}>
                      <div className="flex flex-row gap-x-2"><h3 className="text-neutral-800 text-[14px] font-medium">Group {index + 1}:</h3><div className="text-neutral-800 text-[14px] font-medium">{group.services.map(service => service.charAt(0).toUpperCase() + service.slice(1)).join(", ")}</div></div>
                      <div className="flex flex-row gap-x-2"><div className="text-neutral-500 text-sm font-normal">Nodes: </div><div className="text-neutral-500 text-sm font-medium">{group.nodes}</div></div>
                      <div className="flex flex-row gap-x-2"><div className="text-neutral-500 text-sm font-normal">CPU: </div><div className="text-neutral-500 text-sm font-medium">{group.estimated_cpu} vCPUs</div></div>
                      <div className="flex flex-row gap-x-2"><div className="text-neutral-500 text-sm font-normal">RAM: </div><div className="text-neutral-500 text-sm font-medium">{group.estimated_ram} GB</div></div>
                      <div className="flex flex-row gap-x-2"><div className="text-neutral-500 text-sm font-normal">Disk Type: </div><div className="text-neutral-500 text-sm font-medium">{group.disk_type}</div></div>
                      <div className="flex flex-row gap-x-2"><div className="text-neutral-500 text-sm font-normal">Disk Size: </div><div className="text-neutral-500 text-sm font-medium">{group.estimated_disk} GB</div></div>
                      <div className="flex flex-row gap-x-2"><div className="text-neutral-500 text-sm font-normal">IOPS: </div><div className="text-neutral-500 text-sm font-medium">{group.estimated_disk_io}</div></div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </motion.div>
        )}
      </div>

      <div className="pl-4 pb-4">
        <button
          onClick={() => {
            if (serviceGroups.length > 0) {
              alert(`Creating cluster with the following configuration:\n${JSON.stringify(serviceGroups, null, 2)}`);
            }
          }}
          className={`w-[184px] h-[48px] rounded-[4px] text-[#4C4C4C] text-[16px] flex items-center justify-center
          ${serviceGroups.length === 0 ? "bg-[#E5E5E5] cursor-not-allowed" : "bg-blue-600 text-white rounded-md hover:bg-blue-700 cursor-pointer"}`}>
          Add this configuration
        </button>
      </div>
    </div>
  );
};

export default RecommendedConfig;
