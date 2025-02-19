import React from "react";

interface ServiceGroup {
  service_group_type: string;
  estimated_ram: string;
  estimated_cpu: string;
  estimated_disk: string;
  estimated_disk_io: string;
}

interface RecommendedConfigProps {
  serviceGroups: ServiceGroup[];
}

const RecommendedConfig: React.FC<RecommendedConfigProps> = ({ serviceGroups }) => {
  return (
    <div className="w-[303px] h-full bg-[#F7FAFF] border border-[#E5E5E5] rounded-[12px] flex flex-col">
      <div className="flex items-center border-b border-[#DFDFDF] pl-4 pr-4 pb-1 mt-[5px]">
        <h2 className="text-[14px] text-[#181D27] h-[28px] flex-1 flex items-center font-medium">Recommended Configuration</h2>
      </div>

      <div className="flex-1 mt-[11px] pl-4 pr-4">
        {serviceGroups.length === 0 ? (
            <p className="text-[#535862] text-[14px] w-[257px]">
            Start entering the information to get your recommendation here.
            </p>
        ) : (
          <ul className="space-y-2 text-[14px]">
            {serviceGroups.map((group, index) => (
              <li key={index} className="p-2">
                <p className="font-semibold text-gray-700">{group.service_group_type}</p>
                <p className="text-gray-600">RAM: {group.estimated_ram}</p>
                <p className="text-gray-600">CPU: {group.estimated_cpu}</p>
                <p className="text-gray-600">Disk: {group.estimated_disk}</p>
                <p className="text-gray-600">Disk I/O: {group.estimated_disk_io}</p>
              </li>
            ))}
          </ul>
        )}
      </div>

      <div className="pl-4 pb-4">
        <button className={`w-[184px] h-[48px] rounded-[4px] text-[#4C4C4C] text-[16px] flex items-center justify-center cursor-pointer
          ${serviceGroups.length === 0 ? "bg-[#E5E5E5]" : "bg-blue-600 text-white rounded-md hover:bg-blue-700"}`}>
          Add this configuration
        </button>
      </div>
    </div>
  );
};

export default RecommendedConfig;
