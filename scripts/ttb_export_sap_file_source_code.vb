
Option Explicit

Public m_rst As New ADODB.Recordset

Public ComName As String
Public sAddress As String
Public DebAcc As String
Public dateEff As String
Public schedule As String
Public Product As String
Public TaxID As String
Public myFile As String

Dim NeedSplitItems As Boolean
Dim Msg, Title, WkBookNM, WkSheetNM, ATSFile, MediaFile, CompanyCode, CompanyName As String
Dim BranchCode, AccountNo, EffectiveDate, ServiceType As String
Dim ATS As Range
Dim FirstSeq As Long
Dim NumberOfRows As Integer
Dim TXNLine As String
Dim WHTLIne As String

'
'
'Edit แก้ไขโดย สุขิต ม.ค. 10 และ ม.ค. 12
'Date 07/06/2012
'

Sub ProtectNo()                      'Unprotectsheet
   ActualPwd = "SMARTEXCEL"          'this password use nowhere
   If ActiveSheet.ProtectContents Then ActiveSheet.Unprotect Password:=ActualPwd
End Sub
Sub ProtectYes()                     'ProtectSheet
   ActiveSheet.Protect Password:=ActualPwd, DrawingObjects:=True, contents:=True, Scenarios:=True
End Sub
Sub Abort(ByVal M As String)
  ' Title = " Payment "
  ' Response = MsgBox(M, vbCritical, Title)
   'End
End Sub
Sub PreRequisiteFields()
  If Cells(2, 5) = Empty Then Call Abort("File name is empty") Else MediaFile = Cells(2, 5)
  If Cells(2, 2) = Empty Then Call Abort("Company name is empty") Else CompanyName = Cells(2, 2)
  If Cells(2, 4) = Empty Then Call Abort("Effective date is empty") Else EffectiveDate = Cells(2, 4)
End Sub
Sub ValidateChq()
Dim r As Integer
    r = 4
    Do
      If Cells(r, 24) = 0 Then
      Call Abort("ช่องทางการเรียกเก็บ Charge On row ·เลข: " & r & "")
      End If
      
      If Left(Cells(r, 29), 3) = "MCP" And Cells(r, 18) = 0 Then
      Call Abort("¡ไม่พบข้อมูล กรุณากรอกข้อมูล Row ที่: " & r & "")
      End If

      If Left(Cells(r, 29), 3) = "MCP" And Cells(r, 19) = 0 Then
      Call Abort("ไม่ได้ระบุ Delivery Method Row ที่: " & r & "")
      End If

      If Left(Cells(r, 29), 3) = "MCP" And Cells(r, 20) = 0 Then
      Call Abort("ช่องทาง Pick Up Br. Row ·เลข: " & r & "")
      End If
     
      r = r + 1
    Loop Until Cells(r, 2) = Empty
End Sub
Sub ValidateSmart()
Dim r As Integer
    r = 4
    Do

    Select Case (Left(Cells(r, 16), 3))
      Case "002", "004", "006", "011", "014", "017", "018", "022", "024", "025", "065", "069", "070", "071", "073"     'A/C = 10 digits
        If (Left(Cells(r, 29), 3) = "MCL") And (Len(Cells(r, 15)) <> 10) Then
        Call Abort("à�?¢·�?èº�?­ª�?¼�?é�?�?ºà§�?¹ row ·�?è: " & r & " ä�?èà·è�?¡�?º 10 �?�?�?¡, try again")
        End If
      Case "067"                                                                               'A/C = 14 digits
        If (Left(Cells(r, 29), 3) = "MCL") And (Len(Cells(r, 15)) <> 14) Then
        Call Abort("à�?¢·�?èº�?­ª�?¼�?é�?�?ºà§�?¹ row ·�?è: " & r & " ä�?èà·è�?¡�?º 14 �?�?�?¡, try again")
        End If
      Case "030", "052"                                                                               'A/C = 15 digits
        If (Left(Cells(r, 29), 3) = "MCL") And (Len(Cells(r, 15)) <> 15) Then
        Call Abort("à�?¢·�?èº�?­ª�?¼�?é�?�?ºà§�?¹ row ·�?è: " & r & " ä�?èà·è�?¡�?º 15 �?�?�?¡, try again")
        End If
      Case "034", "031", "033"                                                                 'A/C = 12 digits
        If (Left(Cells(r, 29), 3) = "MCL") And (Len(Cells(r, 15)) <> 12) Then
        Call Abort("à�?¢·�?èº�?­ª�?¼�?é�?�?ºà§�?¹ row ·�?è: " & r & " ä�?èà·è�?¡�?º 12 �?�?�?¡, try again")
        End If
      Case "027"                                                                         'A/C = 8 digits
        If (Left(Cells(r, 29), 3) = "MCL") And (Len(Cells(r, 15)) <> 8) Then
        Call Abort("à�?¢·�?èº�?­ª�?¼�?é�?�?ºà§�?¹ row ·�?è: " & r & " ä�?èà·è�?¡�?º 8 �?�?�?¡, try again")
        End If
      'Case "005"                                                                                'A/C = 6 digits
       ' If (Left(Cells(r, 29), 3) = "MCL") And (Len(Cells(r, 15)) <> 6) Then
        'Call Abort("à�?¢·�?èº�?­ª�?¼�?é�?�?ºà§�?¹ row ·�?è: " & r & " ä�?èà·è�?¡�?º 6 �?�?�?¡, try again")
        'End If

      Case Else
      End Select
      
      r = r + 1
    Loop Until Cells(r, 1) = Empty
End Sub
Sub ValidateBNT()
Dim r As Integer
    r = 4
    Do
    Select Case (Left(Cells(r, 16), 3))
      Case "002", "004", "006", "011", "014", "017", "018", "022", "024", "025", "065", "069", "070", "071", "073"   'A/C = 10 digits
        If (Left(Cells(r, 29), 3) = "BNT") And (Len(Cells(r, 15)) <> 10) Then
        Call Abort("à�?¢·�?èº�?­ª�?¼�?é�?�?ºà§�?¹ row ·�?è: " & r & " ä�?èà·è�?¡�?º 10 �?�?�?¡, try again")
        End If
      Case "067"                                                                         'A/C = 14 digits
        If (Left(Cells(r, 29), 3) = "BNT") And (Len(Cells(r, 15)) <> 14) Then
        Call Abort("à�?¢·�?èº�?­ª�?¼�?é�?�?ºà§�?¹ row ·�?è: " & r & " ä�?èà·è�?¡�?º 14 �?�?�?¡, try again")
        End If
      Case "030", "052"                                                                               'A/C = 15 digits
        If (Left(Cells(r, 29), 3) = "BNT") And (Len(Cells(r, 15)) <> 15) Then
        Call Abort("à�?¢·�?èº�?­ª�?¼�?é�?�?ºà§�?¹ row ·�?è: " & r & " ä�?èà·è�?¡�?º 15 �?�?�?¡, try again")
        End If
      Case "034", "031", "033"                                                                  'A/C = 12 digits
        If (Left(Cells(r, 29), 3) = "BNT") And (Len(Cells(r, 15)) <> 12) Then
        Call Abort("à�?¢·�?èº�?­ª�?¼�?é�?�?ºà§�?¹ row ·�?è: " & r & " ä�?èà·è�?¡�?º 12 �?�?�?¡, try again")
        End If
      Case "027"                                                                         'A/C = 8 digits
        If (Left(Cells(r, 29), 3) = "BNT") And (Len(Cells(r, 15)) <> 8) Then
        Call Abort("à�?¢·�?èº�?­ª�?¼�?é�?�?ºà§�?¹ row ·�?è: " & r & " ä�?èà·è�?¡�?º 8 �?�?�?¡, try again")
        End If
     ' Case "005"                                                                                'A/C = 6 digits
       ' If (Left(Cells(r, 29), 3) = "BNT") And (Len(Cells(r, 15)) <> 6) Then
      '  Call Abort("à�?¢·�?èº�?­ª�?¼�?é�?�?ºà§�?¹ row ·�?è: " & r & " ä�?èà·è�?¡�?º 6 �?�?�?¡, try again")
     '   End If

      Case Else
      End Select
      
      r = r + 1
    Loop Until Cells(r, 1) = Empty
End Sub

Sub ValidateSCS()
Dim r As Integer
    r = 4
    Do
      
      If Left(Cells(r, 29), 3) = "SCS" And Cells(r, 30) = "" Then
      Call Abort("¡�?�?³�?�?�?º�?�?�?º�?�?�?�?�?º Smart Sameday Row ·�?è: " & r & "")
      End If
     
      r = r + 1
    Loop Until Cells(r, 2) = Empty
End Sub

Function Blank(ByVal n As Integer) As String
Dim i As Integer
Dim TStr, TempStr As String
   TStr = ""
   For i = 1 To n: TStr = TStr & " ": Next i
   Blank = TStr
End Function

'Function ReceiveBranchCode(ByVal BnkCode, BrchCode As String) As String
'Dim TempStr As String
'   TempStr = ""
'   If (BnkCode = "002") Or (BnkCode = "004") Or (BnkCode = "006") Or (BnkCode = "011") Or (BnkCode = "014") Or _
'      (BnkCode = "015") Or (BnkCode = "020") Or (BnkCode = "022") Or (BnkCode = "024") Or (BnkCode = "025") Or _
'      (BnkCode = "033") Or (BnkCode = "034") Or (BnkCode = "065") Or (BnkCode = "069") Or (BnkCode = "073") Then
'      TempStr = "0" & BrchCode
'   ElseIf (BnkCode = "030") Or (BnkCode = "067") Or (BnkCode = "072") Then
'        TempStr = BrchCode
'     Else
'        TempStr = "0001"
'   End If
'   ReceiveBranchCode = TempStr
'End Function

'Function ReceiveAcCode(ByVal BnkCode, AcCode As String) As String
'Dim TempStr As String
'   TempStr = ""
'   If (BnkCode = "002") Or (BnkCode = "004") Or (BnkCode = "006") Or (BnkCode = "011") Or (BnkCode = "014") Or _
'      (BnkCode = "015") Or (BnkCode = "022") Or (BnkCode = "024") Or (BnkCode = "025") Or _
'      (BnkCode = "034") Or (BnkCode = "065") Or (BnkCode = "069") Or (BnkCode = "073") Then
'      TempStr = "0" & AcCode
'  End If
'  If BnkCode = "020" Then TempStr = AcCode                      'Standard Chartered Thai
'  If BnkCode = "030" Then TempStr = Right(AcCode, 11)           'Gov Saving Bank
'  If BnkCode = "033" Then TempStr = "00" & Right(AcCode, 9)     'Gov Housing Bank
'  If BnkCode = "005" Then TempStr = "00000" & AcCode            'ABN AMRO
'  If BnkCode = "028" Then                                       'CALYON
'     If Right(AcCode, 5) = "19000" Then
'         TempStr = Mid(AcCode, 3, 5) & "119000"
'     Else
'         TempStr = Mid(AcCode, 3, 5) & "150000"
'     End If
'  End If
'  If BnkCode = "008" Then TempStr = "0" & AcCode                'JP morgan CHASE
'  If BnkCode = "017" Then TempStr = "0" & AcCode                'CITIBank
'  If BnkCode = "018" Then TempStr = "000" & AcCode              'Sumitomo
'  If BnkCode = "039" Then TempStr = AcCode                      'Mitsuho
'  If BnkCode = "010" Then TempStr = AcCode                      'Tokyo Mitsubishi
'  If BnkCode = "007" Then TempStr = AcCode                      'Standard Chartered
'  If BnkCode = "027" Then TempStr = "000" & AcCode              'Bank of America
'  If BnkCode = "031" Then TempStr = Mid(AcCode, 2, 11)          'HSBC
'  If BnkCode = "032" Then TempStr = "0" & AcCode                'Deutch bank
'  If BnkCode = "026" Then TempStr = "0" & AcCode                'China Bank
'  If BnkCode = "067" Then TempStr = "0" & Right(AcCode, 10)     'Tisco Bank
'  If BnkCode = "072" Then TempStr = "0" & Right(AcCode, 10)     'GE MONEY RETAIL
'  ReceiveAcCode = TempStr
'End Function

Function AmountStr(ByVal Amt As Double, ByVal n As Integer) As String   'insert 0 infront of amt
Dim i, l, k As Integer
Dim TStr, TempStr As String
'Edit by sukit
'20/02/2009 For round up 2 digits only
   TempStr = CStr(Round(Amt, 2))
   l = Len(TempStr)
   On Error Resume Next
   k = InStr(1, TempStr, ".")
   If k > 0 Then
     If l - k = 1 Then
        TempStr = TempStr & "0"
     End If
     TempStr = Left(TempStr, k - 1) & "." & Right(TempStr, 2)
   Else
    TempStr = TempStr & "." & "00"
   End If
   l = Len(TempStr)
   TStr = ""
   For i = 1 To n - l: TStr = TStr & "0": Next i
   AmountStr = TStr & TempStr
End Function

Function NameStr(ByVal NM As String, ByVal n As Integer) As String
Dim i, l, k As Integer
Dim TStr, TempStr As String
   k = InStr(1, NM, ",")
   If k > 0 Then TempStr = Left(NM, k - 1) & Right(NM, Len(NM) - k) Else TempStr = NM
   l = Len(TempStr)
   If l >= n Then
      NameStr = Left(TempStr, n)
  Else
      NameStr = TempStr & Blank(n - l)
  End If
End Function

Sub SaveDataTextFile()
Dim r, c, l As Integer
Dim MediaFileNM, recBRCode As String
   MediaFileNM = Workbooks(WkBookNM).Worksheets(WkSheetNM).Cells(2, 5)
   Cells(1, 1) = CStr(MediaHeaderRecord)
   r = 4                    'row of worksheets(wksheetnm)
   Do
    With Workbooks(WkBookNM).Worksheets(WkSheetNM)
         If Left(.Cells(r, 16), 3) = "030" Then                  'GOV Saving Bank
            recBRCode = ReceiveBranchCode("030", Left(.Cells(r, 15), 4))
                ElseIf Left(.Cells(r, 16), 3) = "067" Then       'Tisco
                  recBRCode = ReceiveBranchCode("067", Left(.Cells(r, 15), 4))
                ElseIf Left(.Cells(r, 16), 3) = "033" Then   'GHB
                  recBRCode = ReceiveBranchCode("0330", Left(.Cells(r, 15), 3))
                   ' ElseIf Left(.Cells(r, 16), 3) = "072" Then   'GE Money
                    'recBRCode = ReceiveBranchCode("072", Left(.Cells(r, 15), 4))
                        Else
                        recBRCode = ReceiveBranchCode(Left(.Cells(r, 16), 3), Left(.Cells(r, 15), 3))
         End If

    TXNLine = "TXN" & NameStr(.Cells(2, 2), 120) & _
                NameStr(.Cells(r, 2), 130) & NameStr(.Cells(r, 17), 40) & _
                NameStr(.Cells(r, 18), 170) & NameStr(.Cells(r, 4), 16) & _
                Left(.Cells(2, 4), 2) & Mid(.Cells(2, 4), 4, 2) & Right(.Cells(2, 4), 4) & _
                Left(.Cells(2, 4), 2) & Mid(.Cells(2, 4), 4, 2) & Right(.Cells(2, 4), 4) & _
                "THB" & Blank(50) & "0" & NameStr(.Cells(2, 3), 19) & _
                AmountStr(CDbl(.Cells(r, 3)), 15) & NameStr(Left(.Cells(r, 16), 3), 3) & _
                recBRCode & Blank(9) & NameStr(ReceiveAcCode(Left(.Cells(r, 16), 3), .Cells(r, 15)), 20) & _
                "04" & "00" & NameStr(.Cells(r, 19), 2) & NameStr(.Cells(r, 20), 20) & _
                NameStr(.Cells(r, 21), 5) & NameStr(.Cells(r, 22), 50) & NameStr(.Cells(r, 23), 50) & NameStr(.Cells(r, 24), 50) & _
                NameStr(.Cells(r, 25), 13) & NameStr(.Cells(r, 30), 3) & NameStr(.Cells(r, 31), 5) & Blank(34) & NameStr(.Cells(r, 5), 105) & _
                Blank(295) & "END"
    WHTLIne = "WHT" & NameStr(.Cells(r, 26), 13) & NameStr(.Cells(2, 1), 13) & NameStr(.Cells(r, 27), 2) & _
              AmountStr(CDbl(.Cells(r, 7)), 15) & Blank(2) & NameStr(.Cells(r, 8), 35) & _
              NameStr(.Cells(r, 9), 5) & AmountStr(CDbl(.Cells(r, 10)), 15) & _
              AmountStr(CDbl(.Cells(r, 11)), 15) & Blank(2) & NameStr(.Cells(r, 12), 35) & _
              NameStr(.Cells(r, 13), 5) & AmountStr(CDbl(.Cells(r, 14)), 15) & _
              Blank(144) & NameStr(.Cells(2, 2), 120) & NameStr(.Cells(2, 6), 160) & _
              NameStr(.Cells(r, 28), 120) & NameStr(.Cells(r, 29), 160) & _
              NameStr(.Cells(r, 6), 20) & Blank(938)
                    
    End With
    
    Cells(2 * r, 1) = TXNLine
    Cells((2 * r) + 1, 1) = WHTLIne
    r = r + 1
   
   Loop Until Workbooks(WkBookNM).Worksheets(WkSheetNM).Cells(r, 2) = Empty
   On Error Resume Next
   Application.DisplayAlerts = False
   ActiveWorkbook.SaveAs Filename:=MediaFileNM, FileFormat:=xlTextMSDOS
   If Err <> 0 Then MsgBox "Save Error! : " & MediaFileNM & " code = " & Err
   ActiveWorkbook.Close 'SaveChanges:=False
End Sub

Sub GenerateTextFile()
Dim act As String

Dim TempBook As String
   If Cells(4, 2) = Empty Then
      Call Abort("ä�?è¾º�?�?�?¡�?�?, ä�?è�?�?�?�?�?¶�?�?é�?§ä¿�?ìä´é")
   Else
       act = Application.GetSaveAsFilename("TMBFile")
       If UCase(act) <> "FALSE" Then
           ComName = Sheet1.Cells(2, 2)
           sAddress = Sheet1.Cells(2, 6)
           TaxID = Sheet1.Cells(2, 1)
           'DebAcc = Left(Sheet1.cmbDebitAcc.Value, 10)
       End If
    
       Call PreRequisiteFields
       Call ValidateChq
       Call ValidateSmart
       Call ValidateBNT
       Call ValidateSCS
       
       Call CreateDataSource
       Call CreateText
       
     '  WkBookNM = ActiveWorkbook.Name
     '  WkSheetNM = ActiveSheet.Name
      ' Application.Workbooks.Add
       'Call SaveDataTextFile
     '  Application.Workbooks(WkBookNM).Activate
     '  Worksheets(WkSheetNM).Activate
     
       Msg = "Create text file:  " & Cells(2, 5) & Chr(10) & _
                   "successfully."
       ' Response = MsgBox(Msg, vbInformation, Title)
       ' Cells(4, 1).Activate
   End If
End Sub
Public Function CreateDataSource() As Boolean
Dim rst As ADODB.Recordset
Set rst = New ADODB.Recordset
Dim i As Integer

On Error Resume Next


rst.Fields.Append "ItemNo", vbString, 10 '0
rst.Fields.Append "Payee", vbString, 50  '1
rst.Fields.Append "ACCNO", vbString, 15  '2
rst.Fields.Append "BANK", vbString, 20   '3
rst.Fields.Append "Amt", vbDouble, 30    '4
rst.Fields.Append "Ref", vbString, 20    '5
rst.Fields.Append "DocReq", vbString, 50 '6
rst.Fields.Append "AdvMode", vbString, 10 '7
rst.Fields.Append "Fax", vbString, 10     ' 8
rst.Fields.Append "Mail", vbString, 10    '9
rst.Fields.Append "SMS", vbString, 10     '10
rst.Fields.Append "ComName", vbString, 50 '11
rst.Fields.Append "Address", vbString, 50  '12
rst.Fields.Append "ComAccNo", vbString, 10 '13
rst.Fields.Append "DateEff", vbDate, 10    '14
rst.Fields.Append "Schedule", vbString, 10 '15
rst.Fields.Append "Product", vbString, 10  '16
rst.Fields.Append "TaxID", vbString, 15    '17
rst.Fields.Append "ChargeOn", vbString, 5    '18
i = 5

 rst.Open

Do While Not Sheet1.Cells(i, 2) = ""

  
   rst.AddNew
   rst.Fields(0).Value = Sheet1.Cells(i, 1)
   rst.Fields(1).Value = Sheet1.Cells(i, 2)
   rst.Fields(2).Value = CheckGSB(Left(Sheet1.Cells(i, 16), 3), Replace(Sheet1.Cells(i, 15), "-", ""))
   rst.Fields(3).Value = Left(Sheet1.Cells(i, 16), 3)
   rst.Fields(4).Value = Sheet1.Cells(i, 3)
   rst.Fields(5).Value = Sheet1.Cells(i, 4)
   rst.Fields(6).Value = ""
   rst.Fields(7).Value = Sheet1.Cells(i, 21)
   rst.Fields(8).Value = Sheet1.Cells(i, 22)
   rst.Fields(9).Value = Sheet1.Cells(i, 23)
   rst.Fields(10).Value = Sheet1.Cells(i, 24)
   rst.Fields(11).Value = ComName
   rst.Fields(12).Value = sAddress
   rst.Fields(13).Value = DebAcc
   rst.Fields(14).Value = dateEff
   rst.Fields(15).Value = schedule
   rst.Fields(16).Value = Product
   rst.Fields(17).Value = TaxID
   rst.Fields(18).Value = Sheet1.Cells(i, 25)
   rst.Update
   rst.MoveNext
   i = i + 1
   
Loop
Set m_rst = rst
If Err <> 0 Then
   CreateDataSource = False
   Else
   CreateDataSource = True
End If
End Function
Public Sub CreateText()
Dim recBRCode As String
Dim fnum As Integer
Dim r As Integer
Dim strData As String
If Dir(myFile, vbNormal) <> "" Then
 Kill (myFile)
End If

fnum = FreeFile
Open myFile For Binary As fnum
r = 2
m_rst.MoveFirst
Do While Not m_rst.EOF
                    recBRCode = ReceiveBRCode(Left(m_rst.Fields("BANK").Value, 3), m_rst.Fields("ACCNO").Value)
                    strData = strData & "TXN"
                    strData = strData & AddBlank(m_rst.Fields("ComName").Value, 120)   ' Payer name
                    strData = strData & AddBlank(m_rst.Fields("Payee").Value, 130)     ' Beneficiary name
                    strData = strData & AddBlank(" ", 40)                                ' mail to name
                    strData = strData & AddBlank(" ", 40)                              'Beneficiary address 1
                    strData = strData & AddBlank(" ", 40)                              'Beneficiary address 2
                    strData = strData & AddBlank(" ", 40)                              'Beneficiary address 3
                    strData = strData & AddBlank(" ", 40)                              'Beneficiary address 4
                    strData = strData & AddBlank(" ", 10)                              'Zip code
                    strData = strData & AddBlank(m_rst.Fields("Ref").Value, 16)       'Customer Ref
                    strData = strData & Format(CDate(m_rst.Fields("DateEff").Value), "ddMMyyyy") 'Date Efftive
                    strData = strData & Format(CDate(m_rst.Fields("DateEff").Value), "ddMMyyyy")  'Date Pick up
                    strData = strData & "THB"                   '
                    strData = strData & AddBlank(" ", 50)    '
                    strData = strData & "0" & AddBlank(m_rst.Fields("ComAccNo").Value, 19)      ' Debit Account
                    strData = strData & AmountStr(m_rst.Fields("Amt").Value, 15)                                 ' Amount
                    strData = strData & AddBlank(Left(m_rst.Fields("BANK").Value, 3), 3)             'Beneficiary Bank
                    strData = strData & recBRCode & AddBlank(" ", 9)
                    strData = strData & AddBlank(ReceiveAcCode(Left(m_rst.Fields("BANK").Value, 3), m_rst.Fields("ACCNO").Value), 20)     ' Ben Account
                    strData = strData & "04" & "00"
                    strData = strData & AddBlank(" ", 2)
                    strData = strData & AddBlank(" ", 20)                                          ' Pickup Location
                    strData = strData & AddBlank(m_rst.Fields("AdvMode").Value, 5)
                    strData = strData & AddBlank(Replace(m_rst.Fields("FAX").Value, "-", ""), 50)
                    strData = strData & AddBlank(Replace(m_rst.Fields("Mail").Value, "", ""), 50)
                    strData = strData & AddBlank(Replace(m_rst.Fields("SMS").Value, "-", ""), 50)
                    strData = strData & AddBlank(UCase(m_rst.Fields("ChargeOn").Value), 13)
                    strData = strData & AddBlank(m_rst.Fields("Product").Value, 3)
                    strData = strData & AddBlank(Left(m_rst.Fields("Schedule").Value, 5), 5)
                    strData = strData & AddBlank(" ", 34)
                    strData = strData & AddBlank(m_rst.Fields("DocReq").Value, 105)
                    strData = strData & AddBlank(" ", 295)
                    strData = strData & "END" & vbCrLf
                    strData = strData & "WHT"
                    strData = strData & AddBlank(" ", 13)                                ' Payee tax
                    strData = strData & AddBlank(m_rst.Fields("TaxID").Value, 13)                                'Tax ID
                    strData = strData & AddBlank(" ", 2)
                    strData = strData & AmountStr(0, 15)
                    strData = strData & AddBlank(" ", 2)
                    strData = strData & AddBlank(" ", 35)
                    strData = strData & AddBlank(" ", 5)
                    strData = strData & AmountStr(0, 15)
                    strData = strData & AmountStr(0, 15)
                    strData = strData & AddBlank(" ", 2)
                    strData = strData & AddBlank(" ", 35)
                    strData = strData & AddBlank(" ", 5)
                    strData = strData & AmountStr(0, 15)
                    strData = strData & AddBlank(" ", 144)
                    strData = strData & AddBlank(m_rst.Fields("ComName").Value, 120)
                    strData = strData & AddBlank(m_rst.Fields("Address").Value, 160)
                    strData = strData & AddBlank(" ", 120)
                    strData = strData & AddBlank(" ", 160)
                    strData = strData & AddBlank(" ", 20)
                    strData = strData & AddBlank(" ", 938) & vbCrLf
              
     m_rst.MoveNext
     r = r + 1
           Loop
           Put fnum, , strData
     Close #fnum
     MsgBox "Create file successfully please see file path : " & vbCrLf & myFile & ""
End Sub
Function CheckGSB(bank As String, acc As String)
If bank = "030" Then
   If Len(acc) = 12 Then
      acc = "999" & acc
   End If
   
End If
CheckGSB = acc
End Function
Function ReceiveBRCode(ByVal BnkCode, BrchCode As String) As String
'Edit by sukit 11/07/2012
'Add iBank

Dim TempStr As String
   TempStr = ""
   If (BnkCode = "002") Or (BnkCode = "004") Or (BnkCode = "006") Or (BnkCode = "011") Or (BnkCode = "014") Or _
      (BnkCode = "015") Or (BnkCode = "020") Or (BnkCode = "022") Or (BnkCode = "024") Or (BnkCode = "025") Or _
      (BnkCode = "033") Or (BnkCode = "065") Or (BnkCode = "066") Or (BnkCode = "071") Or (BnkCode = "073") Then
      TempStr = "0" & Left(BrchCode, 3)
   ElseIf (BnkCode = "030") Or (BnkCode = "067") Or (BnkCode = "072") Or (BnkCode = "069") Or (BnkCode = "052") Then
        TempStr = Left(BrchCode, 4)
    ElseIf (BnkCode = "034") Then
         TempStr = "0000"
    ElseIf (BnkCode = "045") Then
        TempStr = "0010"
     Else
        TempStr = "0001"
   End If
   ReceiveBRCode = TempStr
End Function
Function AddBlank(ByVal sText As String, ByVal iLen As Long) As String
        Dim l As Integer
        Dim i As Integer
        Dim res As String
        res = ""
        sText = Trim(sText)
        l = Len(sText)
        If l > iLen Then
            res = Left(sText, iLen)
        Else
            For i = 1 To iLen
                If Mid(sText, i, 1) = "" Then
                    res = res & " "
                Else
                    res = res & Mid(sText, i, 1)
                End If

            Next
        End If

        AddBlank = res
    End Function
Function ReceiveAcCode(ByVal BnkCode, AcCode As String) As String
Dim TempStr As String
'edit validate - , space in account no

'Edit by sukit 11/07/2012
'Add iBank
'add iBank 066

AcCode = Replace(AcCode, "-", "")
AcCode = Replace(AcCode, " ", "")
   TempStr = ""
   If (BnkCode = "002") Or (BnkCode = "004") Or (BnkCode = "006") Or (BnkCode = "008") Or (BnkCode = "011") Or (BnkCode = "014") Or _
      (BnkCode = "017") Or (BnkCode = "018") Or (BnkCode = "022") Or (BnkCode = "024") Or (BnkCode = "025") Or (BnkCode = "026") Or _
      (BnkCode = "032") Or (BnkCode = "065") Or (BnkCode = "073") Or (BnkCode = "066") Or (BnkCode = "079") Or _
      (BnkCode = "070") Or (BnkCode = "071") Then
      
       TempStr = IIf(Len(AcCode) = 10, "0" & AcCode, AcCode)
  ElseIf BnkCode = "034" Then TempStr = Right(AcCode, 11)               '¸¡�?.
  ElseIf BnkCode = "020" Then TempStr = AcCode                      'Standard Chartered Thai
  ElseIf (BnkCode = "030") Or (BnkCode = "052") Or (BnkCode = "045") Then TempStr = Right(AcCode, 11)          'Gov Saving Bank
  ElseIf BnkCode = "033" Then TempStr = "00" & Right(AcCode, 9)      'Gov Housing Bank
 ' If BnkCode = "005" Then TempStr = "00000" & AcCode            'ABN AMRO
 ' If BnkCode = "028" Then                                       'CALYON
  '   If Right(AcCode, 5) = "19000" Then
   '      TempStr = Mid(AcCode, 3, 5) & "119000"
   '  Else
   '      TempStr = Mid(AcCode, 3, 5) & "150000"
  '   End If
'  End If
  'If BnkCode = "008" Then TempStr = "0" & AcCode                'JP morgan CHASE
 ' If BnkCode = "017" Then TempStr = "0" & AcCode                'CITIBank
 ' If BnkCode = "018" Then TempStr = "0" & AcCode                'Sumitomo
  ElseIf BnkCode = "039" Then TempStr = AcCode                      'Mizuho
 ' If BnkCode = "010" Then TempStr = AcCode                      'Tokyo Mitsubishi
 ' If BnkCode = "007" Then TempStr = AcCode                      'Standard Chartered
  ElseIf BnkCode = "027" Then TempStr = Format(AcCode, "00000000000")  'Bank of America
  ElseIf BnkCode = "031" Then TempStr = Mid(AcCode, 2, 11)          'HSBC
 ' If BnkCode = "032" Then TempStr = "0" & AcCode                'Deutch bank
 ' If BnkCode = "026" Then TempStr = "0" & AcCode                'China Bank
  ElseIf BnkCode = "067" Then TempStr = "0" & Right(AcCode, 10)     'Tisco Bank
  ElseIf BnkCode = "069" Then
    If Len(AcCode) = 14 Then
       TempStr = "0" & Right(AcCode, 10)     'KIATNAKIN
    Else
       TempStr = Format(AcCode, "00000000000")
    End If
  End If
  ReceiveAcCode = TempStr
End Function

 Function CheckDigitTMB(sACC As String, sBank As String) As Boolean
        Dim ai() As String
        Dim bi() As String
        Dim ac As String
        Dim sum As Long
        Dim cdg As Long
        Dim i, j As Integer
        Dim ch As Integer
        Dim sKeyBank As String
        Dim a, b, c, d, e, f, g, h, k, l As Integer
         sACC = Replace(sACC, "-", "")
        sACC = Replace(sACC, " ", "")
        If sBank = "011" Then
        
        If Len(sACC) <= 0 Or Len(sACC) > 10 Then
            CheckDigitTMB = False
            Exit Function
        End If
       

        If Len(sACC) > 10 Or Len(sACC) < 11 Then
           a = Val(Mid(sACC, 1, 1)) * 2
           b = Val(Mid(sACC, 2, 1)) * 1
           c = Val(Mid(sACC, 3, 1)) * 2
           d = Val(Mid(sACC, 4, 1)) * 1
           e = Val(Mid(sACC, 5, 1)) * 2
           f = Val(Mid(sACC, 6, 1)) * 1
           g = Val(Mid(sACC, 7, 1)) * 2
           k = Val(Mid(sACC, 8, 1)) * 1
           l = Val(Mid(sACC, 9, 1)) * 2
           
             a = IIf(Val(a) < 10, Val(a), (Val(a) - 10) + 1)
             b = IIf(Val(b) < 10, Val(b), (Val(b) - 10) + 1)
             c = IIf(Val(c) < 10, Val(c), (Val(c) - 10) + 1)
             d = IIf(Val(d) < 10, Val(d), (Val(d) - 10) + 1)
             e = IIf(Val(e) < 10, Val(e), (Val(e) - 10) + 1)
             f = IIf(Val(f) < 10, Val(f), (Val(f) - 10) + 1)
             g = IIf(Val(g) < 10, Val(g), (Val(g) - 10) + 1)
             k = IIf(Val(k) < 10, Val(k), (Val(k) - 10) + 1)
             l = IIf(Val(l) < 10, Val(l), (Val(l) - 10) + 1)
            
            sum = a + b + c + d + e + f + g + k + l
           
        End If
      
            If sum <= 0 Then
                CheckDigitTMB = False
                Exit Function
            End If
            cdg = Right(sum, 1)
            ch = 10 - cdg
            If Right(sACC, 1) = Right(ch, 1) Then
                CheckDigitTMB = True
            Else
                CheckDigitTMB = False
            End If
            
            Else
            CheckDigitTMB = True
     End If
    End Function
