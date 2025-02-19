import React from 'react';

interface ServiceGroup {
  service_group_type: string;
  estimated_ram: string;
  estimated_cpu: string;
  estimated_disk: string;
  estimated_disk_io: string;
}

interface ServiceGroupsTableProps {
  serviceGroups: ServiceGroup[];
}

const ServiceGroupsTable: React.FC<ServiceGroupsTableProps> = ({ serviceGroups }) => {
  return (
    <div className="mt-8 p-6 bg-white rounded-lg shadow-md">
      <h2 className="text-2xl font-semibold text-gray-700 mb-4">Resource Estimates</h2>
      <table className="min-w-full table-auto border-collapse border border-gray-300 rounded-lg overflow-hidden">
        <thead className="bg-gray-100">
          <tr>
            <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">Service Group</th>
            <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">RAM</th>
            <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">CPU</th>
            <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">Disk</th>
            <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">Disk I/O</th>
          </tr>
        </thead>
        <tbody>
          {serviceGroups.map((group, index) => (
            <tr key={index} className={`text-sm ${index % 2 === 0 ? 'bg-white' : 'bg-gray-50'} hover:bg-gray-100`}>
              <td className="px-6 py-4 border-b border-gray-200">{group.service_group_type}</td>
              <td className="px-6 py-4 border-b border-gray-200">{group.estimated_ram}</td>
              <td className="px-6 py-4 border-b border-gray-200">{group.estimated_cpu}</td>
              <td className="px-6 py-4 border-b border-gray-200">{group.estimated_disk}</td>
              <td className="px-6 py-4 border-b border-gray-200">{group.estimated_disk_io}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default ServiceGroupsTable;
