const express = require('express');
const { PrismaClient } = require('../generated/prisma');
const ExcelJS = require('exceljs');
const path = require('path');
const fs = require('fs');

const router = express.Router();
const prisma = new PrismaClient();

router.get('/export', async (req, res) => {
  try {
    // Set the variable first
    await prisma.$executeRawUnsafe('SET @a := 0');

    // Then run the SELECT query
    const result = await prisma.$queryRawUnsafe(`
      SELECT
        @a:=@a+1 AS No,
        g.category_name AS 'Jenis Project',
        a.project_name AS 'Nama Project',
        a.project_objective AS 'Project Objective',
        d.name AS 'Bisnis Proses',
        b.fullname AS 'Leader',
        i.Member,
        a.project_start AS 'Mulai Project',
        a.project_finish AS 'Selesai Project',
        a.project_created_date AS 'Upload Project',
        f.subholding_name AS 'Sub Holding',
        c.corporate_name AS 'SBU',
        CASE a.project_status
          WHEN "New" THEN "Waiting Approve"
          WHEN "Ongoing" THEN "Approve"
          WHEN "Reject" THEN "Reject"
          WHEN "Done" THEN "Done"
          WHEN "Drop" THEN "Delete"
          ELSE ""
        END AS 'Status Approve',
        h.step_project AS 'Step Project',
        a.project_nqi_potential AS 'NQI Potential',
        a.project_nqi_real AS 'NQI Real',
        CASE a.project_approve_nqi
          WHEN "None" THEN "Belum Diajukan"
          WHEN "Waiting" THEN "Menunggu Persetujuan"
          WHEN "Approved" THEN "Disetujui"
          WHEN "Reject" THEN "Ditolak"
          ELSE ""
        END AS 'Status NQI'
      FROM kaizen_list_project a
      LEFT JOIN users b ON b.username = a.project_pic_nik
      LEFT JOIN corporate c ON b.corporate_id = c.corporate_id
      LEFT JOIN kaizen_ref_bisnis_proses d ON a.project_bispro_id = d.id
      LEFT JOIN departments e ON b.department_id = e.department_id
      LEFT JOIN subholding f ON b.subholding_id = f.subholding_id
      LEFT JOIN kaizen_ref_category g ON a.project_category = g.category_id
      LEFT JOIN (
        SELECT id_project, MAX(step_project) AS step_project
        FROM kaizen_list_file
        GROUP BY id_project
      ) h ON a.project_id = h.id_project
      LEFT JOIN (
        SELECT a.member_project_id, GROUP_CONCAT(b.fullname SEPARATOR ', ') AS Member
        FROM kaizen_list_member a
        JOIN users b ON a.member_id = b.username
        GROUP BY a.member_project_id
      ) i ON a.project_id = i.member_project_id
      WHERE
        b.corporate_id = 'dc6529e8-7c9b-48a4-9ceb-a6df1dae7c78' AND
        DATE_FORMAT(a.project_created_date, '%Y-%m-%d') >= '2025-05-04' AND
        DATE_FORMAT(a.project_created_date, '%Y-%m-%d') <= '2025-05-06'
      ORDER BY a.project_created_date ASC
    `);

    console.log(result);

    // Load template
    const workbook = new ExcelJS.Workbook();
    const templatePath = path.join(__dirname, '..', 'template.xlsx');

    console.log('Loading template from:', templatePath);

    await workbook.xlsx.readFile(templatePath);

    console.log('Template loaded successfully');

    const worksheet = workbook.getWorksheet(1); // asumsi sheet pertama
    const startRow = 4; // baris awal input data (misalnya baris ke-6)

    result.forEach((row, index) => {
      const excelRow = worksheet.getRow(startRow + index);
      excelRow.getCell(1).value = index + 1; // No
      excelRow.getCell(2).value = row['Jenis Project'];
      excelRow.getCell(3).value = row['Nama Project'];
      excelRow.getCell(4).value = row['Project Objective'];
      excelRow.getCell(5).value = row['Bisnis Proses'];
      excelRow.getCell(6).value = row['Leader'];
      excelRow.getCell(7).value = row['Member'];
      excelRow.getCell(8).value = row['Mulai Project'];
      excelRow.getCell(9).value = row['Selesai Project'];
      excelRow.getCell(10).value = row['Upload Project'];
      excelRow.getCell(11).value = row['Sub Holding'];
      excelRow.getCell(12).value = row['SBU'];
      excelRow.getCell(13).value = row['Status Approve'];
      excelRow.getCell(14).value = row['Step Project'];
      excelRow.getCell(15).value = row['NQI Potential'];
      excelRow.getCell(16).value = row['NQI Real'];
      excelRow.getCell(17).value = row['Status NQI'];
      excelRow.commit();
    });

    // Output hasil
    const buffer = await workbook.xlsx.writeBuffer();
    res.setHeader('Content-Type', 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet');
    res.setHeader('Content-Disposition', `attachment; filename=rekap_report.xlsx`);
    res.send(buffer);

  } catch (error) {
    console.error('Error executing raw query:', error);
    res.status(500).json({ error: 'Internal server error' });
  }
});

module.exports = router;
