import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { Card, Table } from "react-bootstrap";
import { SubmitButton } from "../../components/SubmitButton";

import api from "../../service/api";

import { IListBudgets } from "../../types/Budget";
import { IJob } from "../../types/Job";
import { formatDate } from "../../utils/formatDate";
import { formatCurrency } from "../../utils/formatCurrency";
import { formatPhone } from "../../utils/formatPhone";

export default function BudgetDetails() {
  const budgetID = useParams().budgetID;
  const [budget, setBudget] = useState<IListBudgets | null>(null);
  const [budgetJobs, setBudgetJobs] = useState<IJob[]>([]);

  const navigate = useNavigate();
  const goBack = () => navigate("/budgets/list");

  const fetchBudget = async () => {
    try {
      const response = await api.get(`/budgets/${budgetID}`, { withCredentials: true });
      setBudget(response.data);
    } catch (error: any) {
      console.error(error.response.data.error);
    }
  };

  const fetchBudgetJobs = async () => {
    try {
      const response = await api.get(`/budgets/list/jobs/${budgetID}`, { withCredentials: true });
      setBudgetJobs(response.data);
    } catch (error: any) {
      console.error(error.response.data.error);
    }
  };

  useEffect(() => {
    fetchBudget();
    fetchBudgetJobs();
  }, []);

  const handleDeleteBudget = async () => {
    if (budget && confirm("Tem certeza de que deseja excluir o orçamento?")) {
      try {
        const response = await api.delete(`/budgets/${budget.id}`, { withCredentials: true });
        if (response.status === 204) goBack();
      } catch (error: any) {
        alert(error.response.data.error);
      }
    }
  };

  return (
    <Card className="mx-auto mt-4" style={{ width: "24rem" }}>
      <Card.Header className="d-flex justify-content-between align-items-center">
        <span className="material-symbols-outlined" role="button" onClick={goBack}>
          arrow_back
        </span>
        <span className="material-symbols-outlined" role="button">
          download
        </span>
        <span className="material-symbols-outlined" role="button" onClick={handleDeleteBudget}>
          delete
        </span>
      </Card.Header>

      <Card.Body>
        {budget ? (
          <>
            <Card.Title className="mb-3">Orçamento de serviço</Card.Title>

            <p>
              <strong>Data:</strong> {formatDate(budget.budgetDate)} <br />
              <strong>Cliente:</strong> {budget.customerName} <br />
              <strong>Contato:</strong> {formatPhone(budget.customerPhone)} <br />
              <strong>Veículo:</strong> {budget.vehicle} <br />
              <strong>Placa:</strong> {budget.licensePlate} <br />
            </p>

            <Table className="mt-3">
              <thead>
                <tr>
                  <th>Serviço</th>
                  <th>Preço</th>
                </tr>
              </thead>
              <tbody>
                {budgetJobs.map((job, index) => (
                  <tr key={index}>
                    <td>{job.name}</td>
                    <td>{formatCurrency(job.price)}</td>
                  </tr>
                ))}
              </tbody>
            </Table>

            <p>
              <strong>Preço final: {formatCurrency(budget.price)}</strong>
            </p>

            <SubmitButton text="Aprovar orçamento" />
          </>
        ) : (
          <p>Carregando...</p>
        )}
      </Card.Body>
    </Card>
  );
}
