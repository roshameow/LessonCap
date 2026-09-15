"""Generate synthetic slides for a manual LessonCap capture exercise.

Requires reportlab. This is a reference PDF, not a recorded native-app export.
"""
from pathlib import Path
from reportlab.pdfgen import canvas
from reportlab.lib.colors import HexColor

out = Path(__file__).resolve().parents[1] / 'docs/examples'
out.mkdir(parents=True, exist_ok=True)
slides = [
    ('Choose a capture region', 'Keep the slide inside the selected area.',
     ['Open these sample slides in a browser.', 'Use LessonCap to select the slide rectangle.', 'Start capture before moving to the next page.']),
    ('Capture a page change', 'A large content change makes the exercise easy to inspect.',
     ['Advance to this slide after capture starts.', 'Wait several seconds before changing pages again.', 'Observe whether the captured-slide count increases.']),
    ('Stop and export', 'Review the result before relying on it for a real lecture.',
     ['Stop capture in LessonCap.', 'Export the captured pages as a PDF.', 'Check page order, boundaries, duplicates and readability.']),
]
pdf = canvas.Canvas(str(out/'sample-slides.pdf'), pagesize=(1280, 720))
pdf.setTitle('LessonCap - Synthetic Capture Exercise')
pdf.setAuthor('LessonCap project')
for i,(title,subtitle,steps) in enumerate(slides,1):
    pdf.setFillColor(HexColor('#142739'));pdf.rect(0,0,1280,720,fill=1,stroke=0)
    pdf.setFillColor(HexColor('#68d5bd'));pdf.setFont('Helvetica-Bold',14)
    pdf.drawString(70,650,'LESSONCAP  /  SYNTHETIC SAMPLE')
    pdf.setFillColor(HexColor('#ffffff'));pdf.setFont('Helvetica-Bold',44);pdf.drawString(70,555,title)
    pdf.setFillColor(HexColor('#b6cad5'));pdf.setFont('Helvetica',23);pdf.drawString(70,505,subtitle)
    for n,step in enumerate(steps):
        y=390-n*82
        pdf.setFillColor(HexColor('#68d5bd'));pdf.circle(90,y+6,19,fill=1,stroke=0)
        pdf.setFillColor(HexColor('#142739'));pdf.setFont('Helvetica-Bold',17);pdf.drawCentredString(90,y,str(n+1))
        pdf.setFillColor(HexColor('#ffffff'));pdf.setFont('Helvetica',23);pdf.drawString(130,y,step)
    pdf.setStrokeColor(HexColor('#355164'));pdf.line(70,105,1210,105)
    pdf.setFillColor(HexColor('#a7bfcb'));pdf.setFont('Helvetica',14)
    pdf.drawString(70,65,'Reference layout generated with ReportLab; not a native-app capture or accuracy benchmark.')
    pdf.drawRightString(1210,65,f'{i} / {len(slides)}');pdf.showPage()
pdf.save()
print(out/'sample-slides.pdf')
